package util

import (
	"path"
	"path/filepath"
	"strings"
)

func SplitPath(p string, pPrefix string) (string, string, []string, bool) {
	pPrefix = strings.TrimSuffix(pPrefix, "/")
	pSplit := filepath.SplitList(strings.ReplaceAll(strings.Trim(p, "/"), "/", ";"))
	if pPrefix == "" && len(pSplit) > 0 {
		pPrefix = pSplit[0]
	}
	dir, fileName := path.Split(strings.TrimPrefix(p, path.Join("/", pPrefix)))
	dirs := filepath.SplitList(strings.ReplaceAll(strings.Trim(dir, "/"), "/", ";"))
	isRoot := (dir == "" || dir == "/") && fileName == ""
	return dir, fileName, dirs, isRoot
}
