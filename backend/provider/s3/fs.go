package s3

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
)

type FileSystem struct {
}

func (f FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) RemoveAll(ctx context.Context, name string, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) HandleHttp(ctx context.Context, name string, file model.ProviderFile, w *http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
