package provider

import (
	"github.com/czy21/cloud-disk-sync/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem interface {
	Mkdir(ctx context.Context, provider model.Provider, name string, perm os.FileMode) error
	OpenFile(ctx context.Context, provider model.Provider, name string, flag int, perm os.FileMode) (webdav.File, error)
	RemoveAll(ctx context.Context, provider model.Provider, name string) error
	Rename(ctx context.Context, provider model.Provider, oldName, newName string) error
	Stat(ctx context.Context, provider model.Provider, name string) (os.FileInfo, error)
}
