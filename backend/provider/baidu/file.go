package baidu

import (
	"context"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"hash"
	"io"
	"io/fs"
)

type File struct {
	file model.ProviderFile
	ctx  context.Context
}

func (f File) Close() error {
	//TODO implement me
	panic("implement me")
}

func (f File) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Stat() (fs.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
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

func (f File) UploadCreate(md5Hash hash.Hash) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) UploadChunk(fileId string, p []byte, md5Bytes []byte, index int) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) error {
	//TODO implement me
	panic("implement me")
}

func (f File) DownloadCreate() (string, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (m int, err error) {
	//TODO implement me
	panic("implement me")
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	panic("implement me")
}
