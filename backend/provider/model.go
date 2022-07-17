package provider

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem interface {
	Mkdir(ctx context.Context, folder model.ProviderFolderMeta, name string, perm os.FileMode) error
	OpenFile(ctx context.Context, folder model.ProviderFolderMeta, name string, flag int, perm os.FileMode) (webdav.File, error)
	RemoveAll(ctx context.Context, folder model.ProviderFolderMeta, name string) error
	Rename(ctx context.Context, folder model.ProviderFolderMeta, oldName, newName string) error
	Stat(ctx context.Context, folder model.ProviderFolderMeta, name string) (os.FileInfo, error)
}
