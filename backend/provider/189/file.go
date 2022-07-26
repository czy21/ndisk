package _189

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/constant"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
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
	cache.Client.Del(f.Context, cache.GetFileInfoCacheKey(f.Name))
	return nil
}

func (f File) Read(b []byte) (n int, err error) {
	extra := f.Context.Value(constant.HttpExtra).(map[string]interface{})
	fileSize := extra[constant.HttpExtraFileSize].(int64)
	_, _, rangeS, rangeE := util.GetChunk(f.Name, fileSize, int64(len(b)), extra)
	dFunc := func(dUrl string) (int, error) {
		req := http2.GetClient().NewRequest()
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeS, rangeE))
		res, _ := req.Get(dUrl)
		if fileSize == 0 || fileSize == rangeE {
			err = io.EOF
		}
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

func (f File) FileName() string {
	return f.Name
}

func (f File) FileSize() int64 {
	extra := f.Context.Value(constant.HttpExtra).(map[string]interface{})
	fileSize := extra[constant.HttpExtraFileSize].(int64)
	return fileSize
}

func (f File) Create(md5Hash hash.Hash) (fileId string, err error) {
	fileSize := f.FileSize()
	d, fName := path.Split(f.Name)
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, d, f.File)
	var fileMd5 string
	if fileSize == 0 {
		fileMd5 = hex.EncodeToString(md5Hash.Sum(nil))
	}
	res, err := API{}.CreateFile(fileInfo.RemoteName, fName, fileSize, fileMd5)
	return res.UploadFileId, nil
}

func (f File) Commit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) (err error) {
	fileSize := f.FileSize()
	fileMd5 := hex.EncodeToString(md5Hash.Sum(nil))
	sliceMd5 := fileMd5
	if fileSize > int64(chunkLen) {
		sliceMd5 = util.GetMD5Encode(strings.Join(md5s, "\n"))
	}
	err = API{}.CommitFile(fileId, fileSize, fileMd5, sliceMd5)
	return err
}

func (f File) Chunk(fileId string, b []byte, md5Bytes []byte, index int) (n int, err error) {
	chunkLen := len(b)
	err = API{}.UploadChunk(fileId, b, md5Bytes, index+1)
	return chunkLen, err
}

func (f File) Write(b []byte) (n int, err error) {
	panic("implement me")
}

// Uploader upload to remote
type Uploader struct {
	Context context.Context
	File    model.ProviderFile
	io.ReadCloser
}

func (u Uploader) WriteTo(w io.Writer) (n int64, err error) {
	l := limitBuf(u.File.ProviderFolder.Account.PutBuf)
	return util.WriteFull(w, u.ReadCloser, 1024*1024*l)
}

// Downloader download from remote
type Downloader struct {
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	l := limitBuf(d.File.ProviderFolder.Account.GetBuf)
	return util.ReadFull(d.ResponseWriter, r, 1024*1024*l)
}

func limitBuf(val int) int {
	if val < 10 || val > 64 {
		val = 10
	}
	return val
}
