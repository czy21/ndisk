package _189

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	fs1 "io/fs"
	"os"
	"path"
	"strconv"
	"time"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) (err error) {
	api := API{File: file}
	dir, fileName := path.Split(file.Path)
	folder, _ := fs.GetFileInfo(ctx, dir, file)
	err = api.CreateFolder(folder.Id, fileName)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	api := API{File: file}
	_, fName := path.Split(file.Path)
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file)
	err = api.Delete(fileInfo.Id, fName, fileInfo.IsDir)
	return err
}
func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	api := API{}
	oldD, oldFName := path.Split(file.OldName)
	newD, newFName := path.Split(file.Name)
	oldFileInfo, err := fs.GetFileInfo(ctx, file.OldName, file)
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
func (fs FileSystem) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (fileInfo model.FileInfo, err error) {
	remoteName := file.ProviderFolder.RemoteName
	fileInfo = model.FileInfo{Name: name, Id: remoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, err
	}
	if !file.IsRoot {
		api := API{}
		var folder FileListAO
		for _, t := range file.Dirs {
			folder, err = api.GetFolderById(remoteName)
			if err == nil {
				err = fs1.ErrNotExist
			}
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.Id
					err = nil
				}
			}
		}
		if file.FileName != "" {
			folder, err = api.GetFolderById(remoteName)
			if err == nil {
				err = fs1.ErrNotExist
			}
			for _, q := range folder.Files {
				if q.Name == file.FileName {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Size = q.Size
					fileInfo.IsDir = false
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					err = nil
				}
			}
			for _, q := range folder.Folders {
				if q.Name == file.FileName {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					err = nil
				}
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
