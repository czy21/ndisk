package webdav

import (
	"fmt"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/provider"
	"github.com/czy21/cloud-disk-sync/provider/local"
	"github.com/czy21/cloud-disk-sync/web"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"

	"os"
	"strings"
)

type FileSystem struct{}

const localDir = "data"

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	web.LogDav("Mkdir", name)
	p, fs := getProvider(ctx, name)
	return fs.Mkdir(ctx, p, name, perm)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	web.LogDav("OpenFile", name)
	if name == "/" {
		return File{}, nil
	}
	p, fs := getProvider(ctx, name)
	return fs.OpenFile(ctx, p, name, flag, perm)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	web.LogDav("RemoveAll", name)
	p, fs := getProvider(ctx, name)
	return fs.RemoveAll(ctx, p, name)
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	web.LogDav("Rename", fmt.Sprintf("%s => %s", oldName, newName))
	p, fs := getProvider(ctx, newName)
	return fs.Rename(ctx, p, oldName, newName)
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	web.LogDav("Stat", name)
	if name == "/" {
		return FileInfo{isDir: true}, nil
	}
	p, fs := getProvider(ctx, name)
	return fs.Stat(ctx, p, name)
}

func getProvider(ctx context.Context, name string) (model.ProviderFolderMeta, provider.FileSystem) {
	folder := model.ProviderFolderMeta{}
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			folder = t
		}
	}
	if fs := provider.Providers[folder.Account.Kind]; fs != nil {
		return folder, fs
	}
	return folder, local.FileSystem{}
}
