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

func (fs FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) error {
	return webdav.Dir(fs.Dir).Mkdir(ctx, name, perm)
}
func (fs FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) (webdav.File, error) {
	return webdav.Dir(fs.Dir).OpenFile(ctx, name, flag, perm)
}
func (fs FileSystem) RemoveAll(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) error {
	return webdav.Dir(fs.Dir).RemoveAll(ctx, name)
}
func (fs FileSystem) Rename(ctx context.Context, oldName, newName string, folder model.ProviderFolderMeta, oldFilePath string, newFilePath string) error {
	return webdav.Dir(fs.Dir).Rename(ctx, oldName, newName)
}
func (fs FileSystem) Stat(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) (os.FileInfo, error) {
	d := path.Join(fs.Dir, folder.Name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		err = os.MkdirAll(d, os.ModePerm)
		exception.Check(err)
	}
	return webdav.Dir(fs.Dir).Stat(ctx, name)
}
