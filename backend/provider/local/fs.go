package local

import (
	"github.com/czy21/ndisk/exception"
	"github.com/czy21/ndisk/model"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
	"path"
)

type FileSystem struct {
	Dir string
}

func NewFS() FileSystem {
	return FileSystem{Dir: viper.GetString("data.dav")}
}

func (fs FileSystem) Mkdir(ctx context.Context, folder model.ProviderFolderMeta, name string, perm os.FileMode) error {
	return webdav.Dir(fs.Dir).Mkdir(ctx, name, perm)
}
func (fs FileSystem) OpenFile(ctx context.Context, folder model.ProviderFolderMeta, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return webdav.Dir(fs.Dir).OpenFile(ctx, name, flag, perm)
}
func (fs FileSystem) RemoveAll(ctx context.Context, folder model.ProviderFolderMeta, name string) error {
	return webdav.Dir(fs.Dir).RemoveAll(ctx, name)
}
func (fs FileSystem) Rename(ctx context.Context, folder model.ProviderFolderMeta, oldName, newName string) error {
	return webdav.Dir(fs.Dir).Rename(ctx, oldName, newName)
}
func (fs FileSystem) Stat(ctx context.Context, folder model.ProviderFolderMeta, name string) (os.FileInfo, error) {
	d := path.Join(fs.Dir, folder.Name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		err = os.MkdirAll(d, os.ModePerm)
		exception.Check(err)
	}
	return webdav.Dir(fs.Dir).Stat(ctx, name)
}
