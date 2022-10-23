package model

import (
	"context"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
)

type HandlerHttp interface {
	HandleHttp(ctx context.Context, name string, file ProviderFile, w *http.ResponseWriter, r *http.Request)
}

type FileSystem interface {
	Mkdir(ctx context.Context, perm os.FileMode, file ProviderFile) error
	OpenFile(ctx context.Context, flag int, perm os.FileMode, file ProviderFile) (webdav.File, error)
	RemoveAll(ctx context.Context, file ProviderFile) error
	Rename(ctx context.Context, file ProviderFile) error
	Stat(ctx context.Context, file ProviderFile) (os.FileInfo, error)
	GetFileInfo(ctx context.Context, name string, file ProviderFile) (FileInfo, error)
}
