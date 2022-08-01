package webdav

import (
	"fmt"
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider"
	"github.com/czy21/ndisk/provider/local"
	"github.com/czy21/ndisk/util"
	"github.com/czy21/ndisk/web"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"io"
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
func (FileSystem) RemoveAll(ctx context.Context, name string) (err error) {
	web.LogDav("RemoveAll", name)
	f, fs := getProvider(name, "")
	err = fs.RemoveAll(ctx, name, f)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.NewPath))
	return err
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) (err error) {
	web.LogDav("Rename", fmt.Sprintf("src:%s dst:%s", oldName, newName))
	f, fs := getProvider(newName, oldName)
	err = fs.Rename(ctx, oldName, newName, f)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.OldPath))
	return err
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	web.LogDav("Stat", name)
	if name == "/" {
		return model.FileInfoDelegate{FileInfo: model.FileInfo{IsDir: true}}, nil
	}
	p, fs := getProvider(name, "")
	return fs.Stat(ctx, name, p)
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

// Uploader upload to remote
type Uploader struct {
	Context context.Context
	File    model.ProviderFile
	io.ReadCloser
}

func (u Uploader) WriteTo(w io.Writer) (n int64, err error) {
	l := limitBuf(u.File.ProviderFolder.Account.PutBuf)
	return util.WriteFull(w, u.ReadCloser, 1024*1024*l)
}

// Downloader download from remote
type Downloader struct {
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	l := limitBuf(d.File.ProviderFolder.Account.GetBuf)
	return util.ReadFull(d.ResponseWriter, r.(*io.LimitedReader).R, 1024*1024*l)
}

func limitBuf(val int) int {
	if val < 10 || val > 64 {
		val = 10
	}
	return val
}

func HandleHttp(name string, w *http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p, fs := getProvider(name, "")
	if h, ok := fs.(provider.HandlerHttp); ok {
		h.HandleHttp(ctx, name, p, w, r)
		return
	}
	extra := ctx.Value(constant.HttpExtra).(map[string]interface{})
	if extra[constant.HttpExtraMethod] == http.MethodGet {
		(*w).Header().Set("Content-Type", util.GetContentType(p.NewPath))
		*w = Downloader{File: p, ResponseWriter: *w}
	}
	if extra[constant.HttpExtraMethod] == http.MethodPut {
		r.Body = Uploader{File: p, ReadCloser: r.Body}
	}
}
