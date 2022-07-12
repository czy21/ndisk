package webdav

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"os"
)

type CloudFileSystem struct{}

func (CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return errors.New("aaa")
}
func (CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return nil, errors.New("aaa")
}
func (CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	return errors.New("aaa")
}
func (CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return errors.New("aa")
}
func (CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return nil, errors.New("aa")
}
