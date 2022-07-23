package _189

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
)

type File struct {
	Name    string
	File    model.ProviderFile
	Context context.Context
	Extra   map[string]interface{}
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}

func (f File) Close() error {
	return nil
}

func (f File) Read(b []byte) (n int, err error) {
	startIndex := int64(0)
	chunkIndex := int64(len(b))
	if val := f.Extra["chunkIndex"]; val != nil {
		v := val.(int64)
		startIndex = v
		chunkIndex += v
	}
	f.Extra["chunkIndex"] = chunkIndex
	//log.Debugf("%s startIndex: %d chunkIndex: %d", f.Name, startIndex, chunkIndex)
	dFunc := func(dUrl string) (int, error) {
		req := http2.GetClient().NewRequest()
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", startIndex, chunkIndex))
		res, err := req.Get(dUrl)
		return copy(b, res.Body()), err
	}
	if dUrl := f.Extra["dUrl"]; dUrl != nil {
		return dFunc(dUrl.(string))
	}
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	fileInfoVO, err := API{}.GetFileInfoById(fileInfo.RemoteName)
	if !fileInfo.IsDir && fileInfoVO.FileDownloadUrl != "" {
		f.Extra["dUrl"] = fileInfoVO.FileDownloadUrl
		return dFunc(fileInfoVO.FileDownloadUrl)
	}
	log.Error(err)
	return len(b), err
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	return fileInfo.Size, err
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	fileInfo, _ := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	folder, err := API{}.GetFolderById(fileInfo.RemoteName)
	var fileInfos []fs.FileInfo
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, model.FileInfoProxy{
			FileInfo: model.FileInfo{
				Name:  t.Name,
				IsDir: true,
			},
		})
	}
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, model.FileInfoProxy{
			FileInfo: model.FileInfo{
				Name: t.Name,
			},
		})
	}
	return fileInfos, err
}

func chunkUpload(uploadFileId string, md5s []string, md5Sum hash.Hash, data []byte, index int) error {
	var md5Bytes []byte
	dMd5 := md5.New()
	dMd5.Write(data)
	md5Bytes = dMd5.Sum(nil)
	md5Hex := hex.EncodeToString(md5Bytes)
	md5Base64 := base64.StdEncoding.EncodeToString(md5Bytes)
	md5s = append(md5s, strings.ToUpper(md5Hex))
	md5Sum.Write(data)
	var uploadUrlsRes UploadUrlVORes
	err := API{}.UploadRequest("/person/getMultiUploadUrls",
		map[string]string{
			"partInfo":     fmt.Sprintf("%d-%s", index, md5Base64),
			"uploadFileId": uploadFileId,
		}, &uploadUrlsRes)
	if err != nil {
		return nil
	}
	uploadData := uploadUrlsRes.UploadUrls[fmt.Sprintf("partNumber_%d", index)]
	log.Debug(uploadData)
	var uploadHeaders []string
	unscapUploadHeader, _ := url.PathUnescape(uploadData.RequestHeader)
	uploadHeaders = strings.Split(unscapUploadHeader, "&")
	uReq := http2.GetClient().NewRequest().SetBody(bytes.NewReader(data))
	for _, t := range uploadHeaders {
		i := strings.Index(t, "=")
		uReq.Header.Set(t[0:i], t[i+1:])
	}
	uRes, err := uReq.Put(uploadData.RequestURL)
	log.Debug(uRes)
	return err
}

func (f File) Write(p []byte) (n int, err error) {
	var sliceIndex int
	var finishedSize int64
	if f.Extra["sliceIndex"] != nil {
		sliceIndex = f.Extra["sliceIndex"].(int) + 1
	} else {
		sliceIndex = 0
	}
	if f.Extra["finishedSize"] != nil {
		finishedSize = f.Extra["finishedSize"].(int64) + int64(len(p))
	} else {
		finishedSize = int64(len(p))
	}
	f.Extra["sliceIndex"] = sliceIndex
	f.Extra["finishedSize"] = finishedSize
	log.Debug(sliceIndex, len(p), finishedSize)
	//chunkSize := int64(len(p))
	//d, fName := path.Split(f.Name)
	//fileInfo, err := FileSystem{}.GetFileInfo(f.Context, d, f.File)
	//fileSize := f.Context.Value("contentLength").(int64)
	//slices := int(math.Max(1, math.Ceil(float64(chunkSize))/float64(fileSize)))
	//var uploadFileId string
	//md5s := make([]string, 0)
	//md5Sum := md5.New()

	//if f.Extra["uploadFileId"] != nil {
	//	uploadFileId = f.Extra["uploadFileId"].(string)
	//
	//	return len(p), nil
	//}
	//res, err := API{}.CreateUpload(fileInfo.RemoteName, fName, fileSize, chunkSize)
	//if res.UploadFileId != "" {
	//	f.Extra["uploadFileId"] = res.UploadFileId
	//}
	return len(p), nil
}

// Uploader upload to remote
type Uploader struct {
	File model.ProviderFile
	io.ReadCloser
}

func (u Uploader) WriteTo(w io.Writer) (n int64, err error) {
	l := limitBuf(u.File.ProviderFolder.Account.PutBuf)
	return util.CopyN(w, u.ReadCloser, make([]byte, 1024*1024*l))
}

// Downloader download from remote
type Downloader struct {
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	l := limitBuf(d.File.ProviderFolder.Account.GetBuf)
	return util.CopyN(d.ResponseWriter, r, make([]byte, 1024*1024*l))
}

func limitBuf(val int) int {
	if val < 10 || val > 64 {
		val = 10
	}
	return val
}
