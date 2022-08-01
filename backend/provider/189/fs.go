package _189

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	fs1 "io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode, file model.ProviderFile) (err error) {
	d, f := path.Split(file.NewPath)
	folder, _ := fs.GetFileInfo(ctx, d, file)
	err = API{}.CreateFolder(folder.RemoteName, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{Name: name, Context: ctx, File: file}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, name string, file model.ProviderFile) error {
	_, fName := path.Split(file.NewPath)
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	err = API{}.Delete(fileInfo.RemoteName, fName, fileInfo.IsDir)
	return err
}
func (fs FileSystem) Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error {
	oldFileInfo, err := fs.GetFileInfo(ctx, oldName, file)
	_, fName := path.Split(newName)
	if !os.IsNotExist(err) {
		if oldFileInfo.IsDir {
			err = API{}.RenameFolder(oldFileInfo.RemoteName, fName)
		}
		err = API{}.RenameFile(oldFileInfo.RemoteName, fName)
	}
	return err
}
func (fs FileSystem) Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	remoteName := file.ProviderFolder.RemoteName
	fileInfo := model.FileInfo{Name: name, RemoteName: remoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	var err error
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, err
	}
	d, f := path.Split(strings.TrimPrefix(name, path.Join("/", strings.TrimSuffix(file.ProviderFolder.Name, "/"))))
	if d == "" || d == "/" {
		return fileInfo, nil
	}
	d = path.Clean(d)
	ds := strings.Split(strings.Trim(d, "/"), "/")
	api := API{}
	var folder FileListAO
	if d != "/" {
		err = fs1.ErrNotExist
		for _, t := range ds {
			folder, err = api.GetFolderById(remoteName)
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.RemoteName
					err = nil
				}
			}
		}
	}
	if f != "" {
		folder, err = api.GetFolderById(remoteName)
		err = fs1.ErrNotExist
		for _, q := range folder.Files {
			if q.Name == f {
				fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
				fileInfo.Size = q.Size
				fileInfo.IsDir = false
				fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				err = nil
			}
		}
		for _, q := range folder.Folders {
			if q.Name == f {
				fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
				fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
				err = nil
			}
		}
	}

	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return fileInfo, err
}
