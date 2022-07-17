package _189

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, folder model.ProviderFolderMeta, name string, perm os.FileMode) error {
	d, f := path.Split(strings.TrimSuffix(name, "/"))
	parentFolder, _ := getFileInfo(d, folder.RemoteName, folder)
	_, err := API{}.CreateFolder(parentFolder.RemoteName, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, folder model.ProviderFolderMeta, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return File{Name: name, Context: ctx, ProviderFolderMeta: folder}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, folder model.ProviderFolderMeta, name string) error {
	_, f := path.Split(strings.TrimSuffix(name, "/"))
	file, err := getFileInfo(name, folder.RemoteName, folder)
	err = API{}.Delete(file.RemoteName, f, file.IsDir)
	cache.Client.DelPrefix(context.Background(), cache.GetFileInfoCacheKey(name))
	return err
}
func (fs FileSystem) Rename(ctx context.Context, folder model.ProviderFolderMeta, oldName, newName string) error {
	oldResource, err := getFileInfo(oldName, folder.RemoteName, folder)
	_, f := path.Split(newName)
	if !os.IsNotExist(err) {
		err = API{}.RenameFolder(oldResource.RemoteName, f)
	}
	cache.Client.DelPrefix(context.Background(), cache.GetFileInfoCacheKey(oldName))
	return err
}
func (fs FileSystem) Stat(ctx context.Context, folder model.ProviderFolderMeta, name string) (os.FileInfo, error) {
	fileInfo, err := getFileInfo(name, folder.RemoteName, folder)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}
func getFileInfo(name string, remoteName string, folderMeta model.ProviderFolderMeta) (model.FileInfo, error) {
	fileInfo := model.FileInfo{Name: name, RemoteName: remoteName, IsDir: true, ModTime: time.Time(*folderMeta.UpdateTime)}
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
			folder, err = api.getFolderById(remoteName)
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.RemoteName
				}
			}
		}
		if d != "/" && f != "" {
			folder, err = api.getFolderById(remoteName)
			err = fs.ErrNotExist
			for _, q := range folder.Files {
				if q.Name == f {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.Size = q.Size
					fileInfo.IsDir = false
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					err = nil
				}
			}
			for _, q := range folder.Folders {
				if q.Name == f {
					fileInfo.ModTime = time.Time(q.UpdateDate)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					err = nil
				}
			}
		}
	}
	if err == nil {
		cache.Client.SetObj(context.Background(), cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return fileInfo, err
}
