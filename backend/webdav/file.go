package webdav

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type FileSystem struct{}

func (FileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return errors.New("aaa")
}
func (FileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return nil, errors.New("aaa")
}
func (FileSystem) RemoveAll(ctx context.Context, name string) error {
	return errors.New("aaa")
}
func (FileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return errors.New("aa")
}
func (FileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return nil, errors.New("aa")
}
