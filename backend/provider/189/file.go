package _189

import (
	"context"
	"fmt"
	"github.com/czy21/ndisk/model"
	"io"
	"io/fs"
	"strconv"
)

type File struct {
	Name    string
	File    model.ProviderFile
	Context context.Context
	Extra   map[string]interface{}
}

func (f File) Close() error {
	return nil
}

func (f File) Read(p []byte) (n int, err error) {
	return 0, nil
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

func (f File) ReadFrom(r io.Reader) (n int64, err error) {
	//fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	//url, err := API{}.getDownloadFileUrl(fileInfo.RemoteName)
	//req := http.GetClient().NewRequest()
	//res, err := req.Get(url)
	//c, err := r.Read(res.Body())
	data := make([]byte, 100)
	c, err := r.Read(data)
	fmt.Printf(strconv.Itoa(c))
	fmt.Printf(string(data))
	return 0, err
}
