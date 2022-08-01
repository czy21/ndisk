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
	Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error
	OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error)
	RemoveAll(ctx context.Context, file model.ProviderFile) error
	Rename(ctx context.Context, file model.ProviderFile) error
	Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error)
	GetFileInfo(ctx context.Context, name string, file model.ProviderFolderMeta) (model.FileInfo, error)
}
