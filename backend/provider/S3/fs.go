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

func listObjectToChan(client *minio.Client, bucketName string, objectName string, recursive bool, objectsCh chan string) {
	defer close(objectsCh)
	for o := range client.ListObjects(bucketName, objectName, recursive, nil) {
		objectsCh <- o.Key
	}
}

func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	api := API{file}
	fileInfo, err := fs.GetFileInfo(ctx, file.Target.Name, file)
	bucketName := file.ProviderFolder.RemoteName
	objectName := file.Target.RelPath
	client, err := api.GetClient()
	objectsCh := make(chan string)
	if fileInfo.IsDir {
		objectName += "/"
		go listObjectToChan(client, bucketName, objectName, true, objectsCh)
	} else {
		objectsCh <- objectName
	}
	client.RemoveObjects(bucketName, objectsCh)
	return err
}

func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) (err error) {
	api := API{file}
	client, err := api.GetClient()
	bucketName := file.ProviderFolder.RemoteName
	srcPath := file.Source.RelPath
	dstPath := file.Target.RelPath
	srcInfo, err := fs.GetFileInfo(ctx, file.Source.Name, file)
	objectsCh := make(chan string)
	srcDstMap := make(map[string]string)
	if srcInfo.IsDir {
		srcPath += "/"
		dstPath += "/"
		go listObjectToChan(client, bucketName, srcPath, true, objectsCh)
		for t := range objectsCh {
			srcDstMap[t] = strings.ReplaceAll(t, srcPath, dstPath)
		}
	} else {
		srcDstMap[srcPath] = dstPath
	}
	for k, v := range srcDstMap {
		source := minio.NewSourceInfo(bucketName, k, nil)
		target, _ := minio.NewDestinationInfo(bucketName, v, nil, nil)
		err = client.CopyObject(target, source)
		err = client.RemoveObject(bucketName, k)
	}
	return err
}

func (fs FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, file.Target.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	return base.GetFileInfo(ctx, name, file, func(fileInfo *model.FileInfo) error {
		var err error
		fileInfo.Id = path.Join(fileInfo.Id, "/", fileInfo.Rel)
		api := API{file}
		err = statObject(api, file.ProviderFolder.RemoteName, fileInfo.Rel, fileInfo)
		if err != nil {
			err = statObject(api, file.ProviderFolder.RemoteName, path.Join(fileInfo.Rel)+"/", fileInfo)
		}
		if err != nil && err != fs1.ErrNotExist {
			err = fs1.ErrNotExist
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
