package _189

import (
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
	dir, fileName := path.Split(file.Target.Path)
	folder, _ := fs.GetFileInfo(ctx, dir, file)
	err = api.CreateFolder(folder.Id, fileName)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	api := API{File: file}
	_, fName := path.Split(file.Target.Path)
	fileInfo, err := fs.GetFileInfo(ctx, file.Target.Name, file)
	err = api.Delete(fileInfo.Id, fName, fileInfo.IsDir)
	return err
}
func (fs FileSystem) Rename(ctx context.Context, file model.ProviderFile) error {
	api := API{}
	oldD, oldFName := path.Split(file.Source.Name)
	newD, newFName := path.Split(file.Target.Name)
	oldFileInfo, err := fs.GetFileInfo(ctx, file.Source.Name, file)
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
	fileInfo, err := fs.GetFileInfo(ctx, file.Target.Name, file)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	return base.GetFileInfo(ctx, name, file, func(fileInfo *model.FileInfo) error {
		var err error
		remoteName := file.FileInfo.Id
		api := API{}
		var folder FileListAO
		if !file.Target.IsRoot {
			for _, t := range file.Target.DirNames {
				folders, err := api.GetFoldersById(remoteName)
				if err == nil {
					err = fs1.ErrNotExist
				}
				for _, f := range folders {
					if f.Name == t {
						fileInfo.Id = f.Id
						remoteName = fileInfo.Id
						err = nil
					}
				}
			}
			if file.Target.BaseName != "" {
				folder, err = api.GetObjectsById(remoteName, file.Target.BaseName)
				if err == nil {
					err = fs1.ErrNotExist
				}
				for _, q := range folder.Files {
					if q.Name == file.Target.BaseName {
						fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
						fileInfo.Size = q.Size
						fileInfo.IsDir = false
						fileInfo.Id = strconv.FormatInt(q.Id, 10)
						err = nil
					}
				}
				for _, q := range folder.Folders {
					if q.Name == file.Target.BaseName {
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
		return err
	})
}
