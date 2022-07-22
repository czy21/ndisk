package _189

import (
	"context"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	log "github.com/sirupsen/logrus"
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

func (f File) Read(b []byte) (n int, err error) {
	startIndex := int64(0)
	chunkIndex := int64(len(b))
	if val := f.Extra["chunkIndex"]; val != nil {
		v := val.(int64)
		startIndex = v
		chunkIndex += v
	}
	f.Extra["chunkIndex"] = chunkIndex
	//log.Debugf("%s startIndex: %d chunkIndex: %d", f.Name, startIndex, chunkIndex)
	dFunc := func(dUrl string) (int, error) {
		req := http2.GetClient().NewRequest()
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", startIndex, chunkIndex))
		res, err := req.Get(dUrl)
		return copy(b, res.Body()), err
	}
	if dUrl := f.Extra["dUrl"]; dUrl != nil {
		return dFunc(dUrl.(string))
	}
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	fileInfoVO, err := API{}.GetFileInfoById(fileInfo.RemoteName)
	if !fileInfo.IsDir && fileInfoVO.FileDownloadUrl != "" {
		f.Extra["dUrl"] = fileInfoVO.FileDownloadUrl
		return dFunc(fileInfoVO.FileDownloadUrl)
	}
	log.Error(err)
	return len(b), err
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
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	var gf int
	if gf = d.File.ProviderFolder.Account.GetBuf; gf < 8 || gf > 64 {
		gf = 8
	}
	return io.CopyBuffer(d.ResponseWriter, r, make([]byte, 1024*1024*gf))
}
