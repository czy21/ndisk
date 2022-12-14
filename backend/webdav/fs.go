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
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.Target.Path))
	return err
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) (err error) {
	web.LogDav("Rename", fmt.Sprintf("src:%s dst:%s", oldName, newName))
	f, fs := getProvider(newName, oldName)
	err = fs.Rename(ctx, f)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(f.Source.Path))
	return err
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name == "/" {
		return model.FileInfoDelegate{FileInfo: model.FileInfo{IsDir: true}}, nil
	}
	f, fs := getProvider(name, "")
	web.LogDav("Stat", fmt.Sprintf("%s dir:%s fileName:%s dirNames:%s isRoot:%t", name, f.Target.Dir, f.Target.Base, fmt.Sprint(f.Target.Parents), f.Target.IsRoot))
	return fs.Stat(ctx, f)
}

func getProvider(name string, oldName string) (model.ProviderFile, model.FileSystem) {
	if name != "" {
		name = path.Clean(name)
	}
	if oldName != "" {
		oldName = path.Clean(oldName)
	}
	file := model.ProviderFile{
		Target: model.ProviderFileMeta{
			Name: name,
			Path: strings.TrimSuffix(name, "/"),
		},
		Source: model.ProviderFileMeta{
			Name: oldName,
			Path: strings.TrimSuffix(oldName, "/"),
		},
	}
	nameSplit := strings.SplitAfter(file.Target.Name, "/")
	if len(nameSplit) >= 2 {
		for _, t := range providerMetas {
			if path.Join(nameSplit[0:2]...) == path.Join("/", t.Name) {
				file.ProviderFolder = t
			}
		}
	}
	dstDir, dstFileName, dstDirNames, dstRel, dstIsRoot := util.SplitPath(file.Target.Name, path.Join("/", file.ProviderFolder.Name))
	file.Target.Rel = dstRel
	file.Target.Base = dstFileName
	file.Target.Dir = dstDir
	file.Target.Parents = dstDirNames
	file.Target.IsRoot = dstIsRoot

	srcDir, srcFileName, srcDirNames, srcRel, srcIsRoot := util.SplitPath(file.Source.Name, path.Join("/", file.ProviderFolder.Name))
	file.Source.Rel = srcRel
	file.Source.Base = srcFileName
	file.Source.Dir = srcDir
	file.Source.Parents = srcDirNames
	file.Source.IsRoot = srcIsRoot
	if fs := provider.GetProviders()[file.ProviderFolder.Account.Kind]; fs != nil {
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
	if u.ReadCloser == http.NoBody {
		if wr, ok := w.(io.ReaderFrom); ok {
			return wr.ReadFrom(nil)
		}
	}
	return util.Copy(w, u.ReadCloser)
}

// Downloader download from remote
type Downloader struct {
	File model.ProviderFile
	http.ResponseWriter
}

func (d Downloader) ReadFrom(r io.Reader) (n int64, err error) {
	return util.Copy(d.ResponseWriter, r.(*io.LimitedReader).R)
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
		(*w).Header().Set("Content-Type", util.GetContentType(p.Target.Path))
		*w = Downloader{File: p, ResponseWriter: *w}
	case http.MethodPut:
		r.Body = Uploader{File: p, ReadCloser: r.Body}
	}
}
