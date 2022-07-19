package provider

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem interface {
	Mkdir(ctx context.Context, name string, perm os.FileMode, file model.ProviderFile) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error)
	RemoveAll(ctx context.Context, name string, file model.ProviderFile) error
	Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error
	Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error)
}
