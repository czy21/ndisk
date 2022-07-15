package _189

import (
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/util"
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
	env := ctx.Value("env").(map[string]interface{})
	fileInfo := FileInfo{}
	if name == path.Join("/", pctx.Meta.Name) {
		env[name] = FileInfo{isDir: true, remoteName: pctx.Meta.RemoteName}
		fileInfo.isDir = true
		return fileInfo, nil
	}
	d, f := path.Split(name)
	ds := strings.Split(strings.TrimSuffix(strings.TrimPrefix(d, "/"), "/"), "/")[1:]
	if len(ds) == 0 && f == "" {
		env[name] = FileInfo{isDir: true, remoteName: pctx.Meta.RemoteName}
		fileInfo.isDir = true
		return fileInfo, nil
	}
	//cache.Client.Set(context.Background(), "a", 1)
	fileInfo = getFolderId(ds, f, pctx.Meta.RemoteName, API{Client: util.HttpUtil{}.NewClient()})
	env[name] = fileInfo
	return fileInfo, nil
}

func getFolderId(ds []string, fName string, folderId string, api API) FileInfo {
	var folder FileListAO
	fileInfo := iteratorDirs(ds, api, folderId)
	if fName != "" {
		folder = api.queryMeta(fileInfo.remoteName)
		for _, q := range folder.Files {
			if q.Name == fName {
				fileInfo.name = q.Name
				fileInfo.isDir = false
				fileInfo.size = q.Size
				fileInfo.remoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.modTime = time.Time(q.UpdateDate)
			}
		}
		for _, q := range folder.Folders {
			if q.Name == fName {
				fileInfo.name = q.Name
				fileInfo.remoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.isDir = true
				fileInfo.modTime = time.Time(q.UpdateDate)
			}
		}
	}
	return fileInfo
}

func iteratorDirs(ds []string, api API, folderId string) FileInfo {
	var folder FileListAO
	var fileInfo FileInfo
	for _, t := range ds {
		folder = api.queryMeta(folderId)
		for _, q := range folder.Folders {
			if q.Name == t {
				fileInfo.name = q.Name
				fileInfo.isDir = true
				fileInfo.remoteName = strconv.FormatInt(q.Id, 10)
				fileInfo.modTime = time.Time(q.UpdateDate)
				folderId = fileInfo.remoteName
			}
		}
	}
	return fileInfo
}
