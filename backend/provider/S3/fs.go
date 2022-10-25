package S3

import (
	"context"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"github.com/minio/minio-go/v6"
	"golang.org/x/net/webdav"
	fs1 "io/fs"
	"os"
	"path"
	"strings"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	api := API{file}
	client, err := api.GetClient()
	_, err = client.PutObject(file.ProviderFolder.RemoteName, path.Join(file.Target.RelPath)+"/", nil, 0, minio.PutObjectOptions{ContentType: ""})
	return err
}

func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}

func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	api := API{file}
	client, err := api.GetClient()
	err = client.RemoveObjectWithOptions(file.ProviderFolder.RemoteName, file.Target.RelPath, minio.RemoveObjectOptions{})
	return err
}

func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) (err error) {
	src := minio.NewSourceInfo(file.ProviderFolder.RemoteName, file.Source.RelPath, nil)
	dst, err := minio.NewDestinationInfo(file.ProviderFolder.RemoteName, file.Target.RelPath, nil, nil)
	api := API{file}
	client, err := api.GetClient()
	err = client.CopyObject(dst, src)
	err = client.RemoveObjectWithOptions(file.ProviderFolder.RemoteName, file.Source.RelPath, minio.RemoveObjectOptions{})
	return err
}

func (fs FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, file.Target.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	return base.GetFileInfo(ctx, name, file, func(fileInfo *model.FileInfo) error {
		var err error
		fileInfo.Id = path.Join(fileInfo.Id) + strings.ReplaceAll(file.Target.Name, path.Join("/", file.ProviderFolder.Name), "")
		if !file.Target.IsRoot {
			api := API{file}
			err = statObject(api, file.ProviderFolder.RemoteName, file.Target.RelPath, file.FileInfo)
			if err != nil {
				err = statObject(api, file.ProviderFolder.RemoteName, path.Join(file.Target.RelPath)+"/", file.FileInfo)
			}
			return err
		}
		return err
	})
}

func statObject(api API, bucketName string, objectName string, fileInfo *model.FileInfo) error {
	objectInfo, stateErr := api.StatObject(bucketName, objectName)
	if stateErr != nil {
		return stateErr
	}
	err := fs1.ErrNotExist
	if objectInfo.Key != "" {
		fileInfo.ModTime = objectInfo.LastModified
		if !strings.HasSuffix(objectInfo.Key, "/") {
			fileInfo.Size = objectInfo.Size
			fileInfo.IsDir = false
		}
		err = nil
	}
	return err
}
