package _189

import (
	"io/fs"
	"os"
	"time"
)

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

func (c FileInfo) Name() string {
	return c.name
}

func (c FileInfo) Size() int64 {
	return c.size
}

func (c FileInfo) Mode() fs.FileMode {
	return c.mode
}

func (c FileInfo) ModTime() time.Time {
	return c.modTime
}

func (c FileInfo) IsDir() bool {
	return c.isDir
}

func (c FileInfo) Sys() any {
	return c.sys
}
