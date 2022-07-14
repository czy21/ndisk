package webdav

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/provider"
	"github.com/czy21/cloud-disk-sync/repository"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"os"
	"strings"
)

type FileSystem struct{}

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Mkdir: %s", name)
	var p model.ProviderMeta
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].Mkdir(ctx, p, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("OpenFile: %s", name)
	if name == "/" {
		return CloudFile{}, nil
	}
	var p model.ProviderMeta
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].OpenFile(ctx, p, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("RemoveAll: %s", name)
	var p model.ProviderMeta
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].RemoveAll(ctx, p, name)
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	var p model.ProviderMeta
	for _, t := range providerMetas {
		if strings.HasPrefix(newName, "/"+t.Name) {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].Rename(ctx, p, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Stat: %s", name)
	if name == "/" {
		providerMetas = repository.Provider{}.SelectList()
		return FileInfo{isDir: true}, nil
	}
	var p model.ProviderMeta
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			p = t
		}
	}
	return provider.Providers[p.Account.Kind].Stat(ctx, p, name)
}
