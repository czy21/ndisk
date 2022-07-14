package local

import (
	"github.com/czy21/cloud-disk-sync/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
	"path"
)

type FileSystem struct{}

const localDir = "data"

func (FileSystem) Mkdir(ctx context.Context, providerMeta model.ProviderMeta, name string, perm os.FileMode) error {
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, providerMeta model.ProviderMeta, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return webdav.Dir(localDir).OpenFile(ctx, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, providerMeta model.ProviderMeta, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (FileSystem) Rename(ctx context.Context, providerMeta model.ProviderMeta, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, providerMeta model.ProviderMeta, name string) (os.FileInfo, error) {
	d := path.Join(localDir, providerMeta.Name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		_ = os.Mkdir(d, 755)
	}
	return webdav.Dir(localDir).Stat(ctx, name)
}
