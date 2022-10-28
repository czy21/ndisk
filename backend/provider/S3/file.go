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

func (f File) Readdir(count int) (fileInfos []fs.FileInfo, err error) {
	api := API{File: f.File}
	objectInfos, err := api.GetObjects(f.File.ProviderFolder.RemoteName, path.Join(f.File.Target.RelPath)+"/")
	for _, t := range objectInfos {
		if path.Join(t.Key) == path.Join(f.File.Target.RelPath) {
			continue
		}
		fileInfo := model.FileInfo{Name: path.Base(t.Key)}
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
	return api.PutObject(f.File.ProviderFolder.RemoteName, f.File.Target.RelPath, r, f.Size(), minio.PutObjectOptions{ContentType: util.GetContentType(f.Name())})
}

func (f File) WriteTo(w io.Writer) (n int64, err error) {
	api := API{f.File}
	httpMethod := util.GetHttpMethod(f.Ctx)
	objectName := f.File.ProviderFolder.RemoteName
	if httpMethod == "COPY" {
		src := minio.NewSourceInfo(objectName, f.File.Target.RelPath, nil)
		dst, err := minio.NewDestinationInfo(objectName, w.(File).File.Target.RelPath, nil, nil)
		client, err := api.GetClient()
		err = client.CopyObject(dst, src)
		return n, err
	}
	object, err := api.GetObject(objectName, f.File.Target.RelPath, minio.GetObjectOptions{})
	if err != nil {
		return 0, err
	}
	return io.Copy(w, object)
}
