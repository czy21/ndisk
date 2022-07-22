package _189

import (
	"context"
	"fmt"
	"github.com/czy21/ndisk/model"
	"io"
	"io/fs"
	"net/http"
)

type File struct {
	Name    string
	File    model.ProviderFile
	Context context.Context
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {
	return n, nil
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

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}

func (f File) Write(p []byte) (n int, err error) {
	panic("aa")
}

// ReadFrom upload to remote
func (f File) ReadFrom(r io.Reader) (n int64, err error) {
	size := 1024 * 1024 * 10
	buf := make([]byte, size)
	for {
		nr, _ := r.Read(buf)
		if nr <= 0 {
			break
		}
		a := buf[0:nr]
		n += int64(len(a))
	}
	fmt.Println(n)
	return n, err
}

type ResponseWriter struct {
	http.ResponseWriter
}

func (w ResponseWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}
