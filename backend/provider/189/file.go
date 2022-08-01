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
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type File struct {
	file model.ProviderFile
	ctx  context.Context
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.ctx, f.file.Name, f.file.ProviderFolder)
	if f.ctx.Value(constant.HttpExtra).(map[string]interface{})[constant.HttpExtraMethod] == http.MethodPut && os.IsNotExist(err) {
		err = nil
	}
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (f File) Close() error {
	if f.ctx.Value(constant.HttpExtra).(map[string]interface{})[constant.HttpExtraMethod] == http.MethodPut {
		cache.Client.Del(f.ctx, cache.GetFileInfoCacheKey(f.file.Name))
	}
	return nil
}

func (f File) DownloadCreate() (dUrl string, fileSize int64, err error) {
	api := API{File: f.file}
	fileInfo, err := FileSystem{}.GetFileInfo(f.ctx, f.file.Name, f.file.ProviderFolder)
	fileInfoVO, err := api.GetFileInfoById(fileInfo.Id)
	return fileInfoVO.FileDownloadUrl, fileInfoVO.Size, err
}
func (f File) DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (n int, err error) {
	req := http2.GetClient().NewRequest()
	req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	res, _ := req.Get(dUrl)
	return copy(p, res.Body()), err
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.ctx, f.file.Name, f.file.ProviderFolder)
	return fileInfo.Size, err
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	api := API{File: f.file}
	fileInfo, _ := FileSystem{}.GetFileInfo(f.ctx, f.file.Name, f.file.ProviderFolder)
	folder, err := api.GetFolderById(fileInfo.Id)
	var fileInfos []fs.FileInfo
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: model.FileInfo{
				Name:  t.Name,
				IsDir: true,
			},
		})
		fi := model.FileInfo{
			Name:    path.Join(f.file.Name, t.Name),
			ModTime: time.Time(t.UpdateDate).Add(-8 * time.Hour),
			Id:      strconv.FormatInt(t.Id, 10),
			IsDir:   true,
		}
		cache.Client.SetObj(f.ctx, cache.GetFileInfoCacheKey(fi.Name), &fi)
	}
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: model.FileInfo{
				Name: t.Name,
			},
		})
		fi := model.FileInfo{
			Name:    path.Join(f.file.Name, t.Name),
			ModTime: time.Time(t.UpdateDate).Add(-8 * time.Hour),
			Size:    t.Size,
			Id:      strconv.FormatInt(t.Id, 10),
		}
		cache.Client.SetObj(f.ctx, cache.GetFileInfoCacheKey(fi.Name), &fi)
	}
	return fileInfos, err
}

func (f File) Name() string {
	return f.file.Name
}

func (f File) UploadLimitSize() int64 {
	return 1024 * 1024 * 8192
}
func (f File) UploadFileSize() int64 {
	extra := f.ctx.Value(constant.HttpExtra).(map[string]interface{})
	fileSize := extra[constant.HttpExtraFileSize].(int64)
	return fileSize
}

func (f File) UploadCreate(md5Hash hash.Hash) (fileId string, err error) {
	api := API{File: f.file}
	fileSize := f.UploadFileSize()
	d, fName := path.Split(f.file.Name)
	fileInfo, err := FileSystem{}.GetFileInfo(f.ctx, d, f.file.ProviderFolder)
	var fileMd5 string
	if fileSize == 0 {
		fileMd5 = hex.EncodeToString(md5Hash.Sum(nil))
	}
	res, err := api.CreateFile(fileInfo.Id, fName, fileSize, fileMd5)
	return res.UploadFileId, err
}

func (f File) UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) (err error) {
	api := API{File: f.file}
	fileSize := f.UploadFileSize()
	fileMd5 := hex.EncodeToString(md5Hash.Sum(nil))
	sliceMd5 := fileMd5
	if fileSize > int64(chunkLen) {
		sliceMd5 = util.GetMD5Encode(strings.Join(md5s, "\n"))
	}
	err = api.CommitFile(fileId, fileSize, fileMd5, sliceMd5)
	return err
}

func (f File) UploadChunk(fileId string, b []byte, md5Bytes []byte, index int) (n int, err error) {
	api := API{File: f.file}
	chunkLen := len(b)
	err = api.UploadChunk(fileId, b, md5Bytes, index+1)
	return chunkLen, err
}

func (f File) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (f File) Write(b []byte) (n int, err error) {
	panic("implement me")
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	api := API{File: f.file}
	_, srcFName := path.Split(f.file.Name)
	dstD, _ := path.Split(w.(File).Name())
	srcFileInfo, err := FileSystem{}.GetFileInfo(f.ctx, f.file.Name, f.file.ProviderFolder)
	dstFileInfo, err := FileSystem{}.GetFileInfo(f.ctx, dstD, f.file.ProviderFolder)
	err = api.Copy(srcFileInfo.Id, srcFName, srcFileInfo.IsDir, dstFileInfo.Id)
	return srcFileInfo.Size, err
}
