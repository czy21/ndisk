package webdav

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/provider"
	"github.com/czy21/cloud-disk-sync/repository"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"os"
)

type CloudFileSystem struct{}

const localDir = "data"

var folders []model.ProviderFolderDTO

func (CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Mkdir: %s", name)
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("OpenFile: %s", name)
	if name == "/" {
		return CloudFile{}, nil
	}
	return provider.All[ctx.Value("providerKind").(string)].OpenFile(ctx, name, flag, perm)
}
func (CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("RemoveAll: %s", name)
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Stat: %s", name)
	if name == "/" {
		folders = repository.Provider{}.SelectAllForFolder()
		return CloudFileInfo{isDir: true}, nil
	}
	return provider.All[ctx.Value("providerKind").(string)].Stat(ctx, name)
}
