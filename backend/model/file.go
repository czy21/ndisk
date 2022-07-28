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

type FileInfoDelegate struct {
	FileInfo
}

func (c FileInfoDelegate) Name() string {
	return c.FileInfo.Name
}

func (c FileInfoDelegate) Size() int64 {
	return c.FileInfo.Size
}

func (c FileInfoDelegate) Mode() fs.FileMode {
	return c.FileInfo.Mode
}

func (c FileInfoDelegate) ModTime() time.Time {
	return c.FileInfo.ModTime
}

func (c FileInfoDelegate) IsDir() bool {
	return c.FileInfo.IsDir
}

func (c FileInfoDelegate) Sys() any {
	return c.FileInfo.Sys
}

func (c FileInfoDelegate) ContentType(ctx context.Context) (cType string, err error) {
	cType = mime.TypeByExtension(filepath.Ext(c.Name()))
	if cType != "" {
		return cType, err
	}
	return "application/octet-stream", err
}
