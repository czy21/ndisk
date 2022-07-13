package webdav

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
	"io/fs"
	"log"
	"os"
	"time"
)

type CloudFileSystem struct{}
type CloudFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (c CloudFileInfo) Name() string {
	return c.name
}

func (c CloudFileInfo) Size() int64 {
	return c.size
}

func (c CloudFileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c CloudFileInfo) ModTime() time.Time {
	return c.modTime
}

func (c CloudFileInfo) IsDir() bool {
	return c.isDir
}

func (c CloudFileInfo) Sys() any {
	return c.sys
}

func (CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("%s", "Mkdir")
	return errors.New("aaa")
}
func (CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("%s", "OpenFile")
	return nil, errors.New("aaa")
}
func (CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("%s", "RemoveAll")
	return errors.New("aaa")
}
func (CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	return errors.New("aa")
}
func (CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("%s", "Stat")
	return CloudFileInfo{isDir: true}, nil
}
