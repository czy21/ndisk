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
	api := API{File: file}
	d, f := path.Split(file.Path)
	folder, _ := fs.GetFileInfo(ctx, d, file)
	err = api.CreateFolder(folder.Id, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{name: name, ctx: ctx, file: file}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, name string, file model.ProviderFile) error {
	api := API{File: file}
	_, fName := path.Split(file.Path)
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	err = api.Delete(fileInfo.Id, fName, fileInfo.IsDir)
	return err
}
func (fs FileSystem) Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error {
	api := API{}
	oldD, oldFName := path.Split(oldName)
	newD, newFName := path.Split(newName)
	oldFileInfo, err := fs.GetFileInfo(ctx, oldName, file)
	newFoldInfo, err := fs.GetFileInfo(ctx, newD, file)
	if oldD != newD {
		err = api.Move(oldFileInfo.Id, oldFName, oldFileInfo.IsDir, newFoldInfo.Id)
		return err
	}
	if !os.IsNotExist(err) {
		if oldFileInfo.IsDir {
			err = api.RenameFolder(oldFileInfo.Id, newFName)
		} else {
			err = api.RenameFile(oldFileInfo.Id, newFName)
		}
	}
	return err
}
func (fs FileSystem) Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	remoteName := file.ProviderFolder.RemoteName
	fileInfo := model.FileInfo{Name: name, Id: remoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	var err error
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, err
	}
	d, f := path.Split(strings.TrimPrefix(name, path.Join("/", strings.TrimSuffix(file.ProviderFolder.Name, "/"))))
	if (d == "" || d == "/") && f == "" {
		return fileInfo, nil
	}
	d = path.Clean(d)
	ds := strings.Split(strings.Trim(d, "/"), "/")
	api := API{}
	var folder FileListAO
	if d != "/" {
		for _, t := range ds {
			folder, err = api.GetFolderById(remoteName)
			err = fs1.ErrNotExist
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.Id
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
				fileInfo.Id = strconv.FormatInt(q.Id, 10)
				err = nil
			}
		}
		for _, q := range folder.Folders {
			if q.Name == f {
				fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
				fileInfo.Id = strconv.FormatInt(q.Id, 10)
				err = nil
			}
		}
	}
	if os.IsNotExist(err) {
		fileInfo.Id = ""
	}
	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	return fileInfo, err
}
