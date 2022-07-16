package _189

import (
	"context"
	"github.com/czy21/cloud-disk-sync/model"
	"io/fs"

	"time"
)

type FileInfoProxy struct {
	model.FileInfo
}

func (c FileInfoProxy) Name() string {
	return c.FileInfo.Name
}

func (c FileInfoProxy) Size() int64 {
	return c.FileInfo.Size
}

func (c FileInfoProxy) Mode() fs.FileMode {
	return c.FileInfo.Mode
}

func (c FileInfoProxy) ModTime() time.Time {
	return c.FileInfo.ModTime
}

func (c FileInfoProxy) IsDir() bool {
	return c.FileInfo.IsDir
}

func (c FileInfoProxy) Sys() any {
	return c.FileInfo.Sys
}

type File struct {
	Name                   string
	ProviderMetaRemoteName string
	Context                context.Context
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
	api := API{}
	fileInfo, _ := getFileInfo(f.Name, f.ProviderMetaRemoteName)
	folder, err := api.queryMeta(fileInfo.RemoteName)
	var fileInfos []fs.FileInfo
	for _, t := range folder.Files {
		fileInfos = append(fileInfos, FileInfoProxy{
			model.FileInfo{
				Name: t.Name,
				Size: t.Size,
			},
		})
	}
	for _, t := range folder.Folders {
		fileInfos = append(fileInfos, FileInfoProxy{
			model.FileInfo{
				Name:  t.Name,
				IsDir: true,
			},
		})
	}
	return fileInfos, err
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, _ := getFileInfo(f.Name, f.ProviderMetaRemoteName)
	return FileInfoProxy{fileInfo}, nil
}

func (f File) Write(p []byte) (n int, err error) {
	panic("implement me")
}
