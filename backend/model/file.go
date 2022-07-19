package model

import (
	"io/fs"
	"os"
	"time"
)

type FileInfo struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	Mode       os.FileMode  `json:"mode"`
	ModTime    StandardTime `json:"modTime"`
	IsDir      bool         `json:"isDir"`
	Sys        any          `json:"sys"`
	RemoteName string       `json:"remoteName"`
}

type FileInfoProxy struct {
	FileInfo
}

func (c FileInfoProxy) Name() string {
	return c.FileInfo.Name
}

func (c FileInfoProxy) Size() int64 {
	return c.FileInfo.Size
}

func (c FileInfoProxy) Mode() fs.FileMode {
	return c.FileInfo.Mode
}

func (c FileInfoProxy) ModTime() time.Time {
	return time.Time(c.FileInfo.ModTime)
}

func (c FileInfoProxy) IsDir() bool {
	return c.FileInfo.IsDir
}

func (c FileInfoProxy) Sys() any {
	return c.FileInfo.Sys
}
