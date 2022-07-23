package _189

import (
	"github.com/czy21/ndisk/cache"
	"github.com/czy21/ndisk/model"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	fs1 "io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileSystem struct {
}

func (fs FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode, file model.ProviderFile) error {
	d, f := path.Split(file.NewPath)
	folder, _ := fs.GetFileInfo(ctx, d, file)
	_, err := API{}.CreateFolder(folder.RemoteName, f)
	return err
}
func (fs FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{Name: name, Context: ctx, File: file, Extra: make(map[string]interface{})}, nil
}
func (fs FileSystem) RemoveAll(ctx context.Context, name string, file model.ProviderFile) error {
	_, fName := path.Split(file.NewPath)
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	err = API{}.Delete(fileInfo.RemoteName, fName, fileInfo.IsDir)
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(file.NewPath))
	return err
}
func (fs FileSystem) Rename(ctx context.Context, oldName, newName string, file model.ProviderFile) error {
	oldFileInfo, err := fs.GetFileInfo(ctx, oldName, file)
	_, fName := path.Split(newName)
	if !os.IsNotExist(err) {
		err = API{}.RenameFolder(oldFileInfo.RemoteName, fName)
	}
	cache.Client.DelPrefix(ctx, cache.GetFileInfoCacheKey(file.OldPath))
	return err
}
func (fs FileSystem) Stat(ctx context.Context, name string, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, name, file)
	return model.FileInfoProxy{FileInfo: fileInfo}, err
}
func (fs FileSystem) GetFileInfo(ctx context.Context, name string, file model.ProviderFile) (model.FileInfo, error) {
	remoteName := file.ProviderFolder.RemoteName
	fileInfo := model.FileInfo{Name: name, RemoteName: remoteName, IsDir: true, ModTime: *file.ProviderFolder.UpdateTime}
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
			folder, err = api.GetFolderById(remoteName)
			for _, q := range folder.Folders {
				if q.Name == t {
					fileInfo.ModTime = time.Time(q.UpdateDate).Add(-8 * time.Hour)
					fileInfo.RemoteName = strconv.FormatInt(q.Id, 10)
					remoteName = fileInfo.RemoteName
				}
			}
		}
		if d != "/" && f != "" {
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
	}
	if err == nil {
		cache.Client.SetObj(ctx, cache.GetFileInfoCacheKey(name), &fileInfo)
	}
	if err == fs1.ErrNotExist {
		err = nil
	}
	return fileInfo, err
}

func (fs FileSystem) HandleHttp(ctx context.Context, name string, file model.ProviderFile, w *http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		*w = Downloader{File: file, ResponseWriter: *w}
	}
	if r.Method == "PUT" {
		r.Body = Uploader{File: file, ReadCloser: r.Body}
	}
}
