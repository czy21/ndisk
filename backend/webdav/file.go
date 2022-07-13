package webdav

import (
	"io/fs"
	"os"
	"time"
)

type CloudFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (c CloudFileInfo) Name() string {
	return c.name
}

func (c CloudFileInfo) Size() int64 {
	return c.size
}

func (c CloudFileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c CloudFileInfo) ModTime() time.Time {
	return c.modTime
}

func (c CloudFileInfo) IsDir() bool {
	return c.isDir
}

func (c CloudFileInfo) Sys() any {
	return c.sys
}

type CloudFile struct{}

func (f CloudFile) Close() error {
	return nil
}

func (f CloudFile) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f CloudFile) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f CloudFile) Readdir(count int) ([]fs.FileInfo, error) {
	var fileInfos []fs.FileInfo
	for _, t := range ProviderFolders {
		fileInfos = append(fileInfos,
			CloudFileInfo{
				isDir: true,
				name:  t.Name,
			})
	}
	return fileInfos, nil
}

func (f CloudFile) Stat() (fs.FileInfo, error) {
	return CloudFileInfo{isDir: true}, nil
}

func (f CloudFile) Write(p []byte) (n int, err error) {
	panic("implement me")
}
