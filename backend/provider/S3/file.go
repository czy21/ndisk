package S3

import (
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"hash"
	"io"
	"io/fs"
	"path"
	"strings"
)

type File struct {
	base.FileBase
}

func (f File) DownloadCreate() (dUrl string, fileSize int64, err error) {
	panic("")
}
func (f File) DownloadChunk(dUrl string, p []byte, rangeStart int64, rangeEnd int64) (n int, err error) {
	panic("")
}

func (f File) Readdir(count int) (fileInfos []fs.FileInfo, err error) {
	api := API{File: f.File}
	objectInfos, err := api.GetObjects(f.File.ProviderFolder.RemoteName, f.File.Dir)
	for _, t := range objectInfos {
		objectName := path.Base(t.Key)
		id := strings.Join([]string{f.File.FileInfo.Id, t.Key}, "/")
		fileInfo := model.FileInfo{Id: id, Name: objectName}
		if strings.HasSuffix(t.Key, "/") {
			fileInfo.IsDir = true
		} else {
			fileInfo.IsDir = false
			fileInfo.Size = t.Size
			fileInfo.ModTime = t.LastModified
		}
		fileInfos = append(fileInfos, model.FileInfoDelegate{
			FileInfo: fileInfo,
		})
	}
	return fileInfos, err
}

func (f File) UploadCreate(md5Hash hash.Hash) (fileId string, err error) {
	panic("")
}

func (f File) UploadCommit(fileId string, md5Hash hash.Hash, md5s []string, chunkLen int) (err error) {
	panic("")
}

func (f File) UploadChunk(fileId string, b []byte, md5Bytes []byte, index int) (n int, err error) {
	panic("")
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	panic("")
}
