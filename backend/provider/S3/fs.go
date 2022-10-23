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
	//TODO implement me
	account := file.ProviderFolder.Account
	client, _ := minio.New(account.Endpoint, account.UserName, account.Password, false)
	print(client)
	return nil, filepath.SkipDir
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	return model.FileInfo{}, nil
}
