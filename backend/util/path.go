package util

import (
	"path"
	"strings"
)

/*
 @p absoute path
*/
func SplitPath(p string, pPrefix string) (string, string, []string, bool) {
	pSplits := strings.SplitAfter(p, "/")
	if (pPrefix == "" || pPrefix == "/") && p != "" {
		pPrefix = path.Join(pSplits[0:2]...)
	}
	dir, fileName := path.Split(strings.TrimPrefix(p, pPrefix))
	dirNames := strings.Split(strings.Trim(dir, "/"), "/")
	isRoot := (dir == "" || dir == "/") && fileName == ""
	return dir, fileName, dirNames, isRoot
}
