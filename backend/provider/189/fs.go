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
	api := API{Client: util.HttpUtil{}.NewClient()}
	d, f := path.Split(name)
	dirSeq := strings.Split(strings.TrimSuffix(strings.TrimPrefix(d, "/"), "/"), "/")[1:]
	if len(dirSeq) == 0 && f == "" {
		env[name] = FileInfo{isDir: true, remoteName: pctx.Meta.RemoteName}
		fileInfo.isDir = true
		return fileInfo, nil
	}
	remoteName, isDir := getFolderId(dirSeq, f, pctx.Meta.RemoteName, api)
	env[name] = FileInfo{remoteName: remoteName, isDir: isDir}
	fileInfo.isDir = isDir
	return fileInfo, nil
}

func getFolderId(dirSeq []string, fName string, folderId string, api API) (string, bool) {
	var isDir bool
	folder, folderId, isDir := iteratorDirs(dirSeq, api, folderId, isDir)
	if fName != "" {
		folder = api.queryMeta(folderId)
		for _, q := range folder.Files {
			if q.Name == fName {
				folderId = strconv.FormatInt(q.Id, 10)
				isDir = false
			}
		}
		for _, q := range folder.Folders {
			if q.Name == fName {
				folderId = strconv.FormatInt(q.Id, 10)
				isDir = true
			}
		}
	}
	return folderId, isDir
}

func iteratorDirs(dPaths []string, api API, folderId string, isDir bool) (FileListAO, string, bool) {
	var folder FileListAO
	for _, t := range dPaths {
		folder = api.queryMeta(folderId)
		for _, q := range folder.Folders {
			if q.Name == t {
				folderId = strconv.FormatInt(q.Id, 10)
				isDir = true
			}
		}
	}
	return folder, folderId, isDir
}
