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
	dir, fileName := path.Split(file.Target.Name)
	folder, err := fs.GetFileInfo(ctx, dir, file)
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
	srcDir, srcBase := path.Split(file.Source.Name)
	dstDir, dstBase := path.Split(file.Target.Name)
	oldFI, err := fs.GetFileInfo(ctx, file.Source.Name, file)
	newFI, err := fs.GetFileInfo(ctx, dstDir, file)
	if srcDir != dstDir {
		err = api.Move(oldFI.Id, srcBase, oldFI.IsDir, newFI.Id)
		return err
	}
	if !os.IsNotExist(err) {
		if oldFI.IsDir {
			err = api.RenameFolder(oldFI.Id, dstBase)
		} else {
			err = api.RenameFile(oldFI.Id, dstBase)
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
		remoteName := file.ProviderFolder.RemoteName
		var err error
		api := API{}
		for _, t := range fileInfo.Parents {
			folders, aErr := api.GetFoldersById(remoteName)
			if aErr != nil {
				err = aErr
				return err
			}
			for _, f := range folders {
				if f.Name == t {
					fileInfo.Id = f.Id
					remoteName = fileInfo.Id
				}
			}
		}
		if fileInfo.Base != "" {
			err = fs1.ErrNotExist
			folder, aErr := api.GetObjectsById(remoteName, fileInfo.Base)
			if aErr != nil {
				err = aErr
			}
			for _, q := range folder.Files {
				if q.Name == fileInfo.Base {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Size = q.Size
					fileInfo.IsDir = false
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					err = nil
				}
			}
			for _, q := range folder.Folders {
				if q.Name == fileInfo.Base {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.Id = strconv.FormatInt(q.Id, 10)
					err = nil
				}
			}
		}
		if err != nil {
			fileInfo.Id = ""
		}
		return err
	})
}
