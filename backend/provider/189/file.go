package _189

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/util"
	"io/fs"
	"os"
	"time"
)

type FileInfo struct {
	name       string
	size       int64
	mode       os.FileMode
	modTime    time.Time
	isDir      bool
	sys        any
	remoteName string
}

func (c FileInfo) Name() string {
	return c.name
}

func (c FileInfo) Size() int64 {
	return c.size
}

func (c FileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c FileInfo) ModTime() time.Time {
	return c.modTime
}

func (c FileInfo) IsDir() bool {
	return c.isDir
}

func (c FileInfo) Sys() any {
	return c.sys
}

type File struct {
	name string
	pctx model.ProviderContext
	env  map[string]interface{}
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f File) Readdir(count int) ([]fs.FileInfo, error) {
	api := API{Client: util.HttpUtil{}.NewClient()}
	folder, err := api.queryMeta(f.env[f.name].(FileInfo).remoteName)
	var fileInfos []fs.FileInfo
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, FileInfo{
			name:    t.Name,
			size:    t.Size,
			modTime: time.Time(t.UpdateDate),
		})
	}
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, FileInfo{
			name:    t.Name,
			modTime: time.Time(t.UpdateDate),
			isDir:   true,
		})
	}
	return fileInfos, err
}

func (f File) Stat() (fs.FileInfo, error) {
	self := f.env[f.name].(FileInfo)
	return FileInfo{isDir: self.isDir, name: f.name, size: self.size, modTime: self.modTime}, nil
}

func (f File) Write(p []byte) (n int, err error) {
	panic("implement me")
}
