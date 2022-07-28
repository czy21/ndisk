package model

import (
	"context"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name       string      `json:"name"`
	Size       int64       `json:"size"`
	Mode       os.FileMode `json:"mode"`
	ModTime    time.Time   `json:"modTime"`
	IsDir      bool        `json:"isDir"`
	Sys        any         `json:"sys"`
	RemoteName string      `json:"remoteName"`
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
	return c.FileInfo.ModTime
}

func (c FileInfoProxy) IsDir() bool {
	return c.FileInfo.IsDir
}

func (c FileInfoProxy) Sys() any {
	return c.FileInfo.Sys
}

func (c FileInfoProxy) ContentType(ctx context.Context) (string, error) {
	ctype := mime.TypeByExtension(filepath.Ext(c.Name()))
	if ctype != "" {
		return ctype, nil
	}
	return "application/octet-stream", nil
}
