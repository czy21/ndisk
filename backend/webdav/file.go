package webdav

import (
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

const localDir = "data"

func (CloudFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Printf("Mkdir: %s", name)
	return webdav.Dir(localDir).Mkdir(ctx, name, perm)
}
func (CloudFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Printf("OpenFile: %s", name)
	return webdav.Dir(localDir).OpenFile(ctx, name, flag, perm)
}
func (CloudFileSystem) RemoveAll(ctx context.Context, name string) error {
	log.Printf("RemoveAll: %s", name)
	return webdav.Dir(localDir).RemoveAll(ctx, name)
}
func (CloudFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	log.Printf("%s", "Rename")
	return webdav.Dir(localDir).Rename(ctx, oldName, newName)
}
func (CloudFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Printf("Stat: %s", name)
	return webdav.Dir(localDir).Stat(ctx, name)
}
