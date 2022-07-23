package _189

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	http2 "github.com/czy21/ndisk/http"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/util"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"math"
	"net/http"
	"path"
	"strings"
)

type File struct {
	Name    string
	File    model.ProviderFile
	Context context.Context
	Extra   map[string]interface{}
}

func (f File) Stat() (fs.FileInfo, error) {
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, f.Name, f.File)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
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

func (f File) Write(p []byte) (n int, err error) {
	chunkSize := int64(len(p))
	d, fName := path.Split(f.Name)
	fileInfo, err := FileSystem{}.GetFileInfo(f.Context, d, f.File)
	fileSize := f.Context.Value("ContentLength").(int64)
	slices := int(math.Max(1, math.Ceil(float64(chunkSize))/float64(fileSize)))
	var uploadFileId string
	md5s := make([]string, 0)
	md5Sum := md5.New()
	uploadFn := func() {
		var i int
		for i = 0; i < slices; i++ {
			var md5Bytes []byte
			dMd5 := md5.New()
			dMd5.Write(p)
			md5Bytes = dMd5.Sum(nil)
			md5Hex := hex.EncodeToString(md5Bytes)
			md5Base64 := base64.StdEncoding.EncodeToString(md5Bytes)
			md5s = append(md5s, strings.ToUpper(md5Hex))
			md5Sum.Write(p)
			
		}
	}
	if f.Extra["uploadFileId"] != nil {
		uploadFileId = f.Extra["uploadFileId"].(string)

		return len(p), nil
	}
	//res, err := API{}.CreateUpload(fileInfo.RemoteName, fName, fileSize, chunkSize)
	//if res.UploadFileId != "" {
	//	f.Extra["uploadFileId"] = res.UploadFileId
	//}
	return len(p), nil
}

// Uploader upload to remote
type Uploader struct {
	File model.ProviderFile
	io.ReadCloser
}

func (u Uploader) WriteTo(w io.Writer) (n int64, err error) {
	return util.CopyBuffer(w, u.ReadCloser, 10)
}

// Downloader download from remote
type Downloader struct {
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	return util.CopyBuffer(d.ResponseWriter, r, d.File.ProviderFolder.Account.GetBuf)
}
