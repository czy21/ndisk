package baidu

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem struct {
}

func (f FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFolderMeta) (model.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}
