package base

import (
	"context"
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"io/fs"
	"net/http"
	"os"
)

type FileBase struct {
	Ctx  context.Context
	FS   model.FileSystem
	File model.ProviderFile
}

func (f FileBase) Stat() (fs.FileInfo, error) {
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Name, f.File.ProviderFolder)
	if f.Ctx.Value(constant.HttpExtra).(map[string]interface{})[constant.HttpExtraMethod] == http.MethodPut && os.IsNotExist(err) {
		err = nil
	}
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (f FileBase) Close() error {
	if f.Ctx.Value(constant.HttpExtra).(map[string]interface{})[constant.HttpExtraMethod] == http.MethodPut {
		cache.Client.Del(f.Ctx, cache.GetFileInfoCacheKey(f.File.Name))
	}
	return nil
}

func (f FileBase) Seek(offset int64, whence int) (int64, error) {
	fileInfo, err := f.FS.GetFileInfo(f.Ctx, f.File.Name, f.File.ProviderFolder)
	return fileInfo.Size, err
}

func (f FileBase) Name() string {
	return f.File.Name
}

func (f FileBase) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (f FileBase) Write(b []byte) (n int, err error) {
	panic("implement me")
}

func (f FileBase) UploadLimitSize() int64 {
	return 1024 * 1024 * 8192
}

func (f FileBase) UploadFileSize() int64 {
	extra := f.Ctx.Value(constant.HttpExtra).(map[string]interface{})
	return extra[constant.HttpExtraFileSize].(int64)
}
