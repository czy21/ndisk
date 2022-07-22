package _189

import (
	"context"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"io"
	"io/fs"
	"net/http"
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
	chunkSize := len(p)
	var (
		startIndex int64
		endIndex   int64
	)
	if val := f.Extra["downloadSize"]; val != nil {
		v := val.(int64)
		startIndex = v
		v += int64(chunkSize)
		endIndex = v
		f.Extra["downloadSize"] = v
	} else {
		endIndex = int64(chunkSize)
		f.Extra["downloadSize"] = int64(chunkSize)
	}
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	url, err := API{}.GetDownloadFileUrl(fileInfo.RemoteName)
	req := http2.GetClient().NewRequest()
	req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", startIndex, endIndex))
	res, err := req.Get(url)
	return copy(p, res.Body()), nil
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
	size := 1024 * 1024 * 4
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

// Downloader download from remote
type Downloader struct {
	Name string
	File model.ProviderFile
	http.ResponseWriter
	Request *http.Request
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	size := 1024 * 1024 * 4
	buf := make([]byte, size)
	for {
		nr, _ := r.Read(buf)
		if nr <= 0 {
			break
		}
		a := buf[0:nr]
		_, _ = d.Write(a)
		n += int64(len(a))
	}
	return n, err
}
