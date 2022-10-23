package S3

import (
	"context"
	"github.com/czy21/ndisk/model"
	"github.com/minio/minio-go/v6"
	"golang.org/x/net/webdav"
	"os"
	"path/filepath"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	//TODO implement me
	panic("implement me")
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
	fileInfo, _ := fs.GetFileInfo(ctx, file.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, filepath.SkipDir
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (fileiInfo model.FileInfo, err error) {
	err = filepath.SkipDir
	//dir, fileName := path.Split(strings.TrimPrefix(name, path.Join("/", strings.TrimSuffix(file.ProviderFolder.Name, "/"))))
	//dirs := strings.Split(strings.Trim(dir, "/"), "/")
	//
	account := file.ProviderFolder.Account
	client, err := minio.New(account.Endpoint, account.UserName, account.Password, false)
	buckets, err := client.ListBuckets()
	for _, t := range buckets {
		println(t.Name)
	}
	println(client.ListBuckets())
	return model.FileInfo{}, nil
}
