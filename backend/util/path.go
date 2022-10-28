package util

import (
	"path"
	"strings"
)

// SplitPath /*
func SplitPath(p string, pPrefix string) (string, string, []string, string, bool) {
	pSplits := strings.SplitAfter(p, "/")
	if (pPrefix == "" || pPrefix == "/") && p != "" {
		pPrefix = path.Join(pSplits[0:2]...)
	}
	rel := strings.TrimPrefix(p, pPrefix)
	dir, base := path.Split(rel)
	isRoot := (dir == "" || dir == "/") && base == ""
	parents := strings.Split(strings.Trim(dir, "/"), "/")
	if len(parents) == 1 && parents[0] == "" {
		parents = nil
	}
	return dir, base, parents, strings.TrimPrefix(rel, "/"), isRoot
}
