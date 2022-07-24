package webdav

import (
	"fmt"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider"
	"github.com/czy21/ndisk/provider/local"
	"github.com/czy21/ndisk/web"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	fs1 "io/fs"
	"net/http"

	"os"
	"strings"
)

type FileSystem struct{}

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	web.LogDav("Mkdir", name)
	p, fs := getProvider(name, "")
	return fs.Mkdir(ctx, name, perm, p)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	web.LogDav("OpenFile", name)
	if name == "/" {
		return File{Name: name}, nil
	}
	p, fs := getProvider(name, "")
	return fs.OpenFile(ctx, name, flag, perm, p)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	web.LogDav("RemoveAll", name)
	p, fs := getProvider(name, "")
	return fs.RemoveAll(ctx, name, p)
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	web.LogDav("Rename", fmt.Sprintf("src:%s dest:%s", oldName, newName))
	p, fs := getProvider(newName, oldName)
	return fs.Rename(ctx, oldName, newName, p)
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	web.LogDav("Stat", name)
	if name == "/" {
		return model.FileInfoProxy{FileInfo: model.FileInfo{IsDir: true}}, nil
	}
	p, fs := getProvider(name, "")
	f, err := fs.Stat(ctx, name, p)
	if ctx.Value("method") == http.MethodPut {
		if err == fs1.ErrNotExist {
			err = nil
		}
	}
	return f, err
}

func getProvider(name string, oldName string) (model.ProviderFile, provider.FileSystem) {
	file := model.ProviderFile{}
	for _, t := range providerMetas {
		if strings.HasPrefix(name, "/"+t.Name) {
			file.ProviderFolder = t
		}
	}
	file.OldPath = strings.TrimSuffix(oldName, "/")
	file.NewPath = strings.TrimSuffix(name, "/")
	if fs := provider.GetProviders()[file.ProviderFolder.Account.Kind]; fs != nil {
		return file, fs
	}
	return file, local.NewFS()
}

func HandleHttp(name string, w *http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p, fs := getProvider(name, "")
	fs.HandleHttp(ctx, name, p, w, r)

}
