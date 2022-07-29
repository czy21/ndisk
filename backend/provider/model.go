package provider

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
)

type HandlerHttp interface {
	HandleHttp(ctx context.Context, name string, file model.ProviderFile, w *http.ResponseWriter, r *http.Request)
}

type FileSystem interface {
	Mkdir(ctx context.Context, name string, perm os.FileMode, file model.ProviderFile) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error)
	RemoveAll(ctx context.Context, name string, file model.ProviderFile) error
	Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error
	Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error)
	GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error)
}

type API interface {
	SetProvider(file model.ProviderFile)
	Create(parentFileInfo model.FileInfo, name string) error
	Delete(parentFileInfo model.FileInfo, name string) error
	Rename(oldFileInfo model.FileInfo, name string) error
	GetFile(fileInfo model.FileInfo) model.FileInfo
	GetFiles(fileInfo model.FileInfo) []model.FileInfo
}
