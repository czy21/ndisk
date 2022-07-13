package provider

import (
	_189 "github.com/czy21/cloud-disk-sync/provider/189"
	"golang.org/x/net/webdav"
)

var All map[string]webdav.FileSystem

func init() {
	All = make(map[string]webdav.FileSystem)
	All["189"] = _189.FileSystem{}
}
