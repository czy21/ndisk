package base

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystemBase struct {
}

func (fs FileSystemBase) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {

	panic("implement me")
}

func (fs FileSystemBase) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	panic("implement me")
}

func (fs FileSystemBase) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	panic("implement me")
}

func (fs FileSystemBase) Rename(ctx context.Context, file model.ProviderFile) error {
	panic("implement me")
}

func (fs FileSystemBase) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	panic("implement me")
}

func (fs FileSystemBase) GetFileInfo(ctx context.Context, name string, file model.ProviderFolderMeta) (model.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}
