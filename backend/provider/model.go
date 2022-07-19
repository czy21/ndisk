package provider

import (
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem interface {
	Mkdir(ctx context.Context, name string, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) error
	OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) (webdav.File, error)
	RemoveAll(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) error
	Rename(ctx context.Context, oldName, newName string, folder model.ProviderFolderMeta, oldFilePath string, newFilePath string) error
	Stat(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) (os.FileInfo, error)
}
