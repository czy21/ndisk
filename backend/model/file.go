package model

import (
	"os"
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
