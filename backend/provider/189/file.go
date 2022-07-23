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
	"path"
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
	_, _, rangeL, rangeR := util.GetChunk(f.Name, int64(len(b)), f.Extra)
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
	if err != nil {
		return 0, err
	}
	if !fileInfo.IsDir {
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

func uploadChunk(uploadFileId string, md5s []string, md5Sum hash.Hash, data []byte, index int) error {
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
		return err
	}
	uploadData := uploadUrlsRes.UploadUrls[fmt.Sprintf("partNumber_%d", index)]
	log.Debug(uploadData)
	uploadHeader, _ := url.PathUnescape(uploadData.RequestHeader)
	uploadHeaders := strings.Split(uploadHeader, "&")
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
	fileSize := f.Extra[ContentLength].(int64)
	chunkSize := int64(len(p))
	d, fName := path.Split(f.Name)
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, d, f.File)
	_, chunkIndex, _, rangeR := util.GetChunk(f.Name, chunkSize, f.Extra)
	// CreateFile
	var fileId string
	if f.Extra[UploadFileId] != nil {
		fileId = f.Extra[UploadFileId].(string)
	} else {
		res, err := API{}.CreateUpload(fileInfo.RemoteName, fName, fileSize)
		if err != nil {
			return 0, nil
		}
		fileId = res.UploadFileId
		f.Extra[UploadFileId] = res.UploadFileId
	}
	// UploadFile
	var md5Sum hash.Hash
	var md5s []string
	if f.Extra[UploadMd5s] == nil {
		md5s = make([]string, 0)
		md5Sum = md5.New()
		f.Extra[UploadMd5s] = md5s
		f.Extra[UploadMd5Sum] = md5Sum
	} else {
		md5s = f.Extra[UploadMd5s].([]string)
		md5Sum = f.Extra[UploadMd5Sum].(hash.Hash)
	}
	err = uploadChunk(fileId, md5s, md5Sum, p, chunkIndex+1)
	if err != nil {
		return 0, err
	}
	// CommitFile
	if fileSize == rangeR {
		fileMd5 := hex.EncodeToString(md5Sum.Sum(nil))
		err = API{}.CommitFile(fileId, fileMd5, fileMd5)
	}
	return int(chunkSize), err
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
