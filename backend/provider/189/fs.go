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
	return File{name: name, pctx: pctx, env: ctx.Value("env").(map[string]interface{})}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, pctx model.ProviderContext, name string) error {
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (fs FileSystem) Rename(ctx context.Context, pctx model.ProviderContext, oldName, newName string) error {
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (fs FileSystem) Stat(ctx context.Context, pctx model.ProviderContext, name string) (os.FileInfo, error) {
	var err error
	d, f := path.Split(name)
	ds := strings.Split(strings.TrimSuffix(strings.TrimPrefix(d, "/"), "/"), "/")[1:]
	fileInfo := model.FileInfo{Name: name, IsDir: true, RemoteName: pctx.Meta.RemoteName}
	if d == "/" || (len(ds) == 0 && f == "") {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), fileInfo)
		return FileInfoProxy{fileInfo}, nil
	}
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return FileInfoProxy{fileInfo}, nil
	}
	fileInfo, err = getFileInfo(ds, f, pctx.Meta.RemoteName, API{})
	cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), fileInfo)
	return FileInfoProxy{fileInfo}, err
}

func getFileInfo(ds []string, fName string, folderId string, api API) (model.FileInfo, error) {
	var folder FileListAO
	var err error
	fileInfo, err := iteratorDirs(ds, api, folderId)
	if fName != "" {
		folder, err = api.queryMeta(fileInfo.RemoteName)
		for _, q := range folder.Files {
			if q.Name == fName {
				fileInfo.Name = q.Name
				fileInfo.IsDir = false
				fileInfo.Size = q.Size
				fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.ModTime = time.Time(q.UpdateDate)
			}
		}
		for _, q := range folder.Folders {
			if q.Name == fName {
				fileInfo.Name = q.Name
				fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.IsDir = true
				fileInfo.ModTime = time.Time(q.UpdateDate)
			}
		}
	}
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
