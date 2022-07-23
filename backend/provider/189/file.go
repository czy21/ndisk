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
	_, rangeL, rangeR := getChunk(f.Name, int64(len(b)), f.Extra)
	dFunc := func(dUrl string) (int, error) {
		req := http2.GetClient().NewRequest()
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeL, rangeR))
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
	_, _, _ = getChunk(f.Name, int64(len(p)), f.Extra)
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

func getChunk(name string, chunkSize int64, extra map[string]interface{}) (int, int64, int64) {
	chunkI := 0
	rangeL := int64(0)
	rangeR := chunkSize
	if extra["chunkI"] != nil {
		chunkI = extra["chunkI"].(int) + 1
	}
	if val := extra["rangeR"]; val != nil {
		v := val.(int64)
		rangeL = v
		rangeR += v
	}
	extra["chunkI"] = chunkI
	extra["rangeR"] = rangeR
	log.Debugf("%s chunkI: %d chunkS: %d rangeL: %d rangeR: %d", name, chunkI, chunkSize, rangeL, rangeR)
	return chunkI, rangeL, rangeR
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
