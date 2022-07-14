package webdav

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/provider"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"log"
	"os"
	"strings"
)

type FileSystem struct{}

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Mkdir: %s", name)
	p, fs := getProvider(ctx, name)
	return fs.Mkdir(ctx, p, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("OpenFile: %s", name)
	if name == "/" {
		return CloudFile{}, nil
	}
	p, fs := getProvider(ctx, name)
	return fs.OpenFile(ctx, p, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("RemoveAll: %s", name)
	p, fs := getProvider(ctx, name)
	return fs.RemoveAll(ctx, p, name)
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	p, fs := getProvider(ctx, newName)
	return fs.Rename(ctx, p, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Stat: %s", name)
	if name == "/" {
		return FileInfo{isDir: true}, nil
	}
	p, fs := getProvider(ctx, name)
	return fs.Stat(ctx, p, name)
}

func getProvider(ctx context.Context, name string) (model.ProviderContext, provider.FileSystem) {
	p := model.ProviderContext{}
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			p.Meta = t
		}
	}
	return p, provider.Providers[p.Meta.Account.Kind]
}
