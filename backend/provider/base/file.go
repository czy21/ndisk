package base

import (
	"io/fs"
)

type File struct {
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
