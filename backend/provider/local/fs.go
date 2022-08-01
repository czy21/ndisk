package local

import (
	"github.com/czy21/ndisk/exception"
	"github.com/czy21/ndisk/model"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"path"
)

type FileSystem struct {
	Dir string
}

func NewFS() FileSystem {
	return FileSystem{Dir: viper.GetString("data.dav")}
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	return webdav.Dir(fs.Dir).Mkdir(ctx, file.Name, perm)
}
func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return webdav.Dir(fs.Dir).OpenFile(ctx, file.Name, flag, perm)
}
func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	return webdav.Dir(fs.Dir).RemoveAll(ctx, file.Name)
}
func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	return webdav.Dir(fs.Dir).Rename(ctx, file.OldName, file.Name)
}
func (fs FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	d := path.Join(fs.Dir, file.ProviderFolder.Name)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		err = os.MkdirAll(d, os.ModePerm)
		exception.Check(err)
	}
	return webdav.Dir(fs.Dir).Stat(ctx, file.Name)
}

func (fs FileSystem) GetFileInfo(ctx context.Context, name string, providerFile model.ProviderFolderMeta) (model.FileInfo, error) {
	return model.FileInfo{}, nil
}

func (fs FileSystem) HandleHttp(ctx context.Context, name string, file model.ProviderFile, w *http.ResponseWriter, r *http.Request) {

}
