package webdav

import (
	"fmt"
	"github.com/czy21/ndisk/cache"
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
	"path"
	"strings"
)

type FileSystem struct{}

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	web.LogDav("Mkdir", name)
	f, fs := getProvider(name, "")
	return fs.Mkdir(ctx, perm, f)
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	web.LogDav("OpenFile", name)
	if name == "/" {
		return File{Name: name}, nil
	}
	f, fs := getProvider(name, "")
	return fs.OpenFile(ctx, flag, perm, f)
}
func (FileSystem) RemoveAll(ctx context.Context, name string) (err error) {
	web.LogDav("RemoveAll", name)
	f, fs := getProvider(name, "")
	err = fs.RemoveAll(ctx, f)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.Path))
	return err
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) (err error) {
	web.LogDav("Rename", fmt.Sprintf("src:%s dst:%s", oldName, newName))
	f, fs := getProvider(newName, oldName)
	err = fs.Rename(ctx, f)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.OldPath))
	return err
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name == "/" {
		return model.FileInfoDelegate{FileInfo: model.FileInfo{IsDir: true}}, nil
	}
	f, fs := getProvider(name, "")
	web.LogDav("Stat", fmt.Sprintf("%s dir:%s fileName:%s dirNames:%s isRoot:%t", name, f.Dir, f.BaseName, fmt.Sprint(f.DirNames), f.IsRoot))
	return fs.Stat(ctx, f)
}

func getProvider(name string, oldName string) (model.ProviderFile, model.FileSystem) {
	file := model.ProviderFile{
		Name:    name,
		Path:    strings.TrimSuffix(name, "/"),
		OldName: oldName,
		OldPath: strings.TrimSuffix(oldName, "/"),
	}
	rootPath := path.Join(strings.SplitAfter(file.Name, "/")[0:2]...)
	for _, t := range providerMetas {
		if rootPath == path.Join("/", t.Name) {
			file.ProviderFolder = t
		}
	}
	dir, fileName, dirNames, isRoot := util.SplitPath(file.Name, path.Join("/", file.ProviderFolder.Name))
	file.RelPath = strings.TrimPrefix(dir+fileName, "/")
	file.BaseName = fileName
	file.Dir = dir
	file.DirNames = dirNames
	file.IsRoot = isRoot
	if fs := provider.GetProviders()[file.ProviderFolder.Account.Kind]; fs != nil {
		file.FileInfo = &model.FileInfo{Name: name, Id: file.ProviderFolder.RemoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
		return file, fs
	}
	return file, local.NewFS()
}

// Uploader upload to remote
type Uploader struct {
	File model.ProviderFile
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
	if h, ok := fs.(model.HandlerHttp); ok {
		h.HandleHttp(ctx, name, p, w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		(*w).Header().Set("Content-Type", util.GetContentType(p.Path))
		*w = Downloader{File: p, ResponseWriter: *w}
	case http.MethodPut:
		r.Body = Uploader{File: p, ReadCloser: r.Body}
	}
}
