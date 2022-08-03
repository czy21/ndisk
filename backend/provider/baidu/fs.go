package baidu

import (
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/base"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystemBase struct {
}

func (fs FileSystemBase) Mkdir(ctx context.Context, perm os.FileMode, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystemBase) OpenFile(ctx context.Context, flag int, perm os.FileMode, file model.ProviderFile) (webdav.File, error) {
	return File{base.FileBase{Ctx: ctx, File: file, FS: fs}}, nil
}

func (fs FileSystemBase) RemoveAll(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystemBase) Rename(ctx context.Context, file model.ProviderFile) error {
	//TODO implement me
	panic("implement me")
}

func (fs FileSystemBase) Stat(ctx context.Context, file model.ProviderFile) (os.FileInfo, error) {
	fileInfo, err := fs.GetFileInfo(ctx, file.Name, file.ProviderFolder)
	return model.FileInfoDelegate{FileInfo: fileInfo}, err
}

func (fs FileSystemBase) GetFileInfo(ctx context.Context, name string, file model.ProviderFolderMeta) (model.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}
