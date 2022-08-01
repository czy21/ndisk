package util

import (
	"mime"
	"path/filepath"
)

func BoolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func GetContentType(name string) string {
	cType := mime.TypeByExtension(filepath.Ext(name))
	if cType != "" {
		return cType
	}
	return "application/octet-stream"
}
