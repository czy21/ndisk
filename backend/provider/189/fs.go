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

func (fs FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) error {
	d, f := path.Split(strings.TrimSuffix(name, "/"))
	parentFolder, _ := getFileInfo(ctx, d, folder.RemoteName, folder)
	_, err := API{}.CreateFolder(parentFolder.RemoteName, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, folder model.ProviderFolderMeta, filePath string) (webdav.File, error) {
	return File{Name: name, Context: ctx, ProviderFolderMeta: folder}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) error {
	_, fName := path.Split(filePath)
	file, err := getFileInfo(ctx, name, folder.RemoteName, folder)
	err = API{}.Delete(file.RemoteName, fName, file.IsDir)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(filePath))
	return err
}
func (fs FileSystem) Rename(ctx context.Context, oldName, newName string, folder model.ProviderFolderMeta, oldFilePath string, newFilePath string) error {
	oldFileInfo, err := getFileInfo(ctx, oldName, folder.RemoteName, folder)
	_, fName := path.Split(newName)
	if !os.IsNotExist(err) {
		err = API{}.RenameFolder(oldFileInfo.RemoteName, fName)
	}
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(oldFilePath))
	return err
}
func (fs FileSystem) Stat(ctx context.Context, name string, folder model.ProviderFolderMeta, filePath string) (os.FileInfo, error) {
	fileInfo, err := getFileInfo(ctx, name, folder.RemoteName, folder)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}
func getFileInfo(ctx context.Context, name string, remoteName string, folderMeta model.ProviderFolderMeta) (model.FileInfo, error) {
	fileInfo := model.FileInfo{Name: name, RemoteName: remoteName, IsDir: true, ModTime: time.Time(*folderMeta.UpdateTime)}
	var err error
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
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
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return fileInfo, err
}
