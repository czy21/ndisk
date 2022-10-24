package S3

import (
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/czy21/ndisk/util"
	"github.com/minio/minio-go/v6"
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
	objectInfos, err := api.GetObjects(f.File.ProviderFolder.RemoteName, f.File.RelPath)
	for _, t := range objectInfos {
		objectName := path.Base(t.Key)
		id := strings.Join([]string{f.File.FileInfo.Id, f.File.RelPath}, "/")
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

func (f File) ReadFrom(r io.Reader) (n int64, err error) {
	api := API{f.File}
	client, err := api.GetClient()
	objectName := strings.SplitAfterN(f.Name(), "/", 3)[2]
	return client.PutObject(f.File.ProviderFolder.RemoteName, objectName, r, -1, minio.PutObjectOptions{ContentType: util.GetContentType(f.Name())})
}

//WriteTo CopyTo
func (f File) WriteTo(w io.Writer) (n int64, err error) {
	api := API{f.File}
	client, err := api.GetClient()
	objectName := strings.SplitAfterN(f.Name(), "/", 3)[2]
	object, err := client.GetObject(f.File.ProviderFolder.RemoteName, objectName, minio.GetObjectOptions{})
	return io.Copy(w, object)
}
