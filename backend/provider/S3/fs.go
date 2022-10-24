package S3

import (
	"context"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/minio/minio-go/v6"
	"golang.org/x/net/webdav"
	"os"
	"path"
	"strings"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}

func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	return base.GetFileInfo(ctx, name, file, func(fileInfo *model.FileInfo) error {
		fileInfo.Id = strings.Join([]string{path.Join(fileInfo.Id), name}, "")
		var err error
		api := API{file}
		var objectInfos []minio.ObjectInfo
		objectInfos, err = api.GetObjects(file.ProviderFolder.RemoteName, file.Dir)
		for _, t := range objectInfos {
			objectName := path.Base(t.Key)
			if objectName == file.FileName {
				fileInfo.ModTime = t.LastModified
				if strings.HasSuffix(t.Key, "/") {
					fileInfo.IsDir = true
				} else {
					fileInfo.Size = t.Size
					fileInfo.IsDir = false
				}
			}
		}
		return err
	})
}
