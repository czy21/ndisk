package _189

import (
	"github.com/czy21/cloud-disk-sync/cache"
	"github.com/czy21/cloud-disk-sync/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileSystem struct {
}

const localDir = "data"

func (fs FileSystem) Mkdir(ctx context.Context, pctx model.ProviderContext, name string, perm os.FileMode) error {
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (fs FileSystem) OpenFile(ctx context.Context, pctx model.ProviderContext, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return File{Name: name, Context: ctx, ProviderMetaRemoteName: pctx.Meta.RemoteName}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, pctx model.ProviderContext, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (fs FileSystem) Rename(ctx context.Context, pctx model.ProviderContext, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (fs FileSystem) Stat(ctx context.Context, pctx model.ProviderContext, name string) (os.FileInfo, error) {
	fileInfo, err := getFileInfo(name, pctx.Meta.RemoteName)
	return FileInfoProxy{fileInfo}, err
}
func getFileInfo(name string, remoteName string) (model.FileInfo, error) {
	fileInfo := model.FileInfo{Name: name, RemoteName: remoteName, IsDir: true}
	var err error
	if cache.Client.GetObj(context.Background(), cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, err
	}
	d, f := path.Split(name)
	ds := strings.Split(strings.Trim(d, "/"), "/")[1:]
	if d != "/" || len(ds) > 0 || f != "" {
		api := API{}
		var folder FileListAO
		for _, t := range ds {
			folder, err = api.queryMeta(remoteName)
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.RemoteName
				}
			}
		}
		if f != "" {
			folder, err = api.queryMeta(remoteName)
			for _, q := range folder.Files {
				if q.Name == f {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.Size = q.Size
					fileInfo.IsDir = false
				}
			}
			for _, q := range folder.Folders {
				if q.Name == f {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				}
			}
		}
	}
	cache.Client.SetObj(context.Background(), cache.GetFileInfoCacheKey(name), &fileInfo)
	return fileInfo, err
}
