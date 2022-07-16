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
	fileInfo, err := getFileInfo(ctx, name, pctx.Meta.RemoteName)
	return FileInfoProxy{fileInfo}, err
}
func getFileInfo(ctx context.Context, name string, remoteName string) (model.FileInfo, error) {
	var (
		fileInfo model.FileInfo
		err      error
	)
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, nil
	}
	d, f := path.Split(name)
	ds := strings.Split(strings.TrimSuffix(strings.TrimPrefix(d, "/"), "/"), "/")[1:]
	fileInfo = model.FileInfo{Name: name, IsDir: true, RemoteName: remoteName}
	if d != "/" || len(ds) > 0 || f != "" {
		var folder FileListAO
		api := API{}
		fileInfo, err = iteratorDirs(ds, api, remoteName)
		if f != "" {
			folder, err = api.queryMeta(fileInfo.RemoteName)
			for _, q := range folder.Files {
				if q.Name == f {
					fileInfo.Name = q.Name
					fileInfo.IsDir = false
					fileInfo.Size = q.Size
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					fileInfo.ModTime = time.Time(q.UpdateDate)
				}
			}
			for _, q := range folder.Folders {
				if q.Name == f {
					fileInfo.Name = q.Name
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					fileInfo.IsDir = true
					fileInfo.ModTime = time.Time(q.UpdateDate)
				}
			}
		}
	}
	cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), fileInfo)
	return fileInfo, err
}

func iteratorDirs(ds []string, api API, folderId string) (model.FileInfo, error) {
	fileInfo := model.FileInfo{IsDir: true, RemoteName: folderId}
	var (
		folder FileListAO
		err    error
	)
	for _, t := range ds {
		folder, err = api.queryMeta(folderId)
		for _, q := range folder.Folders {
			if q.Name == t {
				fileInfo.Name = q.Name
				fileInfo.IsDir = true
				fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.ModTime = time.Time(q.UpdateDate)
				folderId = fileInfo.RemoteName
			}
		}
	}
	return fileInfo, err
}
