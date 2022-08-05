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
	"strings"
	"time"
)

type FileSystem struct {
	base.FileSystemBase
}

func (fs FileSystem) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) (err error) {
	api := API{File: file}
	d, f := path.Split(file.Path)
	folder, _ := fs.GetFileInfo(ctx, d, file.ProviderFolder)
	err = api.CreateFolder(folder.Id, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	api := API{File: file}
	_, fName := path.Split(file.Path)
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file.ProviderFolder)
	err = api.Delete(fileInfo.Id, fName, fileInfo.IsDir)
	return err
}
func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	api := API{}
	oldD, oldFName := path.Split(file.OldName)
	newD, newFName := path.Split(file.Name)
	oldFileInfo, err := fs.GetFileInfo(ctx, file.OldName, file.ProviderFolder)
	newFoldInfo, err := fs.GetFileInfo(ctx, newD, file.ProviderFolder)
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
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file.ProviderFolder)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, providerFolder model.ProviderFolderMeta) (model.FileInfo, error) {
	remoteName := providerFolder.RemoteName
	fileInfo := model.FileInfo{Name: name, Id: remoteName, IsDir: true, ModTime: *providerFolder.UpdateTime}
	var err error
	if cache.Client.GetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo) {
		return fileInfo, err
	}
	d, f := path.Split(strings.TrimPrefix(name, path.Join("/", strings.TrimSuffix(providerFolder.Name, "/"))))
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
	}
	if f != "" {
		folder, err = api.GetFolderById(remoteName)
		if err == nil {
			err = fs1.ErrNotExist
		}
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
