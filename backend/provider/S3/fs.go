package S3

import (
	"context"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/czy21/ndisk/util"
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

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (fileInfo model.FileInfo, err error) {
	remoteName := file.ProviderFolder.RemoteName
	fileInfo = model.FileInfo{Name: name, Id: remoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	//if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
	//	return fileInfo, err
	//}
	dir, fileName, _, isRoot := util.SplitPath(name, file.ProviderFolder.Name)
	if !isRoot {
		api := API{file}
		var objectInfos []minio.ObjectInfo
		objectInfos, err = api.GetObjects(file.ProviderFolder.RemoteName, dir)
		for _, t := range objectInfos {
			objectName := path.Base(t.Key)
			if objectName == fileName {
				fileInfo.ModTime = t.LastModified
				fileInfo.Id = name
				if strings.HasSuffix(t.Key, "/") {
					fileInfo.IsDir = true
				} else {
					fileInfo.Size = t.Size
					fileInfo.IsDir = false
				}
			}
		}
	}
	return fileInfo, err
}
