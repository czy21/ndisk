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

type FileSystem struct{}

var providers []model.Provider

const localDir = "data"

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Mkdir: %s", name)
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("OpenFile: %s", name)
	if name == "/" {
		return CloudFile{}, nil
	}
	var p model.Provider
	for _, t := range providers {
		if name == "/"+t.Name {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].OpenFile(ctx, p, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("RemoveAll: %s", name)
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Stat: %s", name)
	if name == "/" {
		providers = repository.Provider{}.SelectList()
		return FileInfo{isDir: true}, nil
	}
	var p model.Provider
	for _, t := range providers {
		if name == "/"+t.Name {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].Stat(ctx, p, name)
}
