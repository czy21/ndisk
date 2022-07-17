package local

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
	"path"
)

type FileSystem struct {
}

const localDir = "data/local"

func (FileSystem) Mkdir(ctx context.Context, folder model.ProviderFolderMeta, name string, perm os.FileMode) error {
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, folder model.ProviderFolderMeta, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return webdav.Dir(localDir).OpenFile(ctx, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, folder model.ProviderFolderMeta, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (FileSystem) Rename(ctx context.Context, folder model.ProviderFolderMeta, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, folder model.ProviderFolderMeta, name string) (os.FileInfo, error) {
	d := path.Join(localDir, folder.Name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		_ = os.Mkdir(d, 755)
	}
	return webdav.Dir(localDir).Stat(ctx, name)
}
