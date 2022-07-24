package _189

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/czy21/ndisk/constant"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type File struct {
	Name    string
	File    model.ProviderFile
	Context context.Context
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}

func (f File) Close() error {
	return nil
}

func (f File) Read(b []byte) (n int, err error) {
	extra := f.Context.Value(constant.HttpExtra).(map[string]interface{})
	fileSize := extra[constant.HttpExtraFileSize].(int64)
	_, _, rangeL, rangeR := util.GetChunk(f.Name, fileSize, int64(len(b)), extra)
	dFunc := func(dUrl string) (int, error) {
		req := http2.GetClient().NewRequest()
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeL, rangeR))
		res, err := req.Get(dUrl)
		return copy(b, res.Body()), err
	}
	if dUrl := extra[constant.HttpExtraDownloadUrl]; dUrl != nil {
		return dFunc(dUrl.(string))
	}
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	fileInfoVO, err := API{}.GetFileInfoById(fileInfo.RemoteName)
	if err != nil {
		return 0, err
	}
	if !fileInfo.IsDir {
		extra[constant.HttpExtraDownloadUrl] = fileInfoVO.FileDownloadUrl
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

func (f File) Write(b []byte) (n int, err error) {
	api := API{}
	extra := f.Context.Value(constant.HttpExtra).(map[string]interface{})
	fileSize := extra[constant.HttpExtraFileSize].(int64)
	chunkLen := int64(len(b))
	d, fName := path.Split(f.Name)
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, d, f.File)
	_, chunkIndex, _, rangeR := util.GetChunk(f.Name, fileSize, chunkLen, extra)
	var md5s []string
	var md5Sum hash.Hash
	if extra[constant.HttpExtraMd5s] != nil {
		md5s = extra[constant.HttpExtraMd5s].([]string)
		md5Sum = extra[constant.HttpExtraMd5Sum].(hash.Hash)
	} else {
		md5s = make([]string, 0)
		md5Sum = md5.New()
	}
	// CreateFile
	var fileId string
	if extra[constant.HttpExtraFileId] != nil {
		fileId = extra[constant.HttpExtraFileId].(string)
	} else {
		var fileMd5 string
		if fileSize == 0 {
			md5Sum.Write(b)
			fileMd5 = hex.EncodeToString(md5Sum.Sum(nil))
		}
		res, err := api.CreateUpload(fileInfo.RemoteName, fName, fileSize, fileMd5)
		if err != nil {
			return 0, nil
		}
		fileId = res.UploadFileId
		extra[constant.HttpExtraFileId] = fileId
	}
	chunkMd5Bytes, err := api.UploadChunk(fileId, b, chunkIndex+1)
	if err != nil {
		return 0, err
	}
	md5Hex := hex.EncodeToString(chunkMd5Bytes)
	md5s = append(md5s, strings.ToUpper(md5Hex))
	md5Sum.Write(b)
	extra[constant.HttpExtraMd5s] = md5s
	extra[constant.HttpExtraMd5Sum] = md5Sum

	// CommitFile
	if fileSize == rangeR {
		fileMd5 := hex.EncodeToString(md5Sum.Sum(nil))
		sliceMd5 := fileMd5
		if fileSize > chunkLen {
			sliceMd5 = util.GetMD5Encode(strings.Join(md5s, "\n"))
		}
		err = api.CommitFile(fileId, fileSize, fileMd5, sliceMd5)
		return int(chunkLen), io.EOF
	}
	return int(chunkLen), err
}

// Uploader upload to remote
type Uploader struct {
	Context context.Context
	File    model.ProviderFile
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
