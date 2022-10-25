package util

import (
	"fmt"
	"testing"
)

func PrintSplitPath(p string, pPrefix string) {
	dir, fileName, dirs, isRoot := SplitPath(p, pPrefix)
	fmt.Println(fmt.Sprintf("dir:%s fileName:%s dirNames:%s isRoot:%t", dir, fileName, dirs, isRoot))
}

func TestSplitPath(t *testing.T) {
	p1 := "a/b/c/d"
	p2 := "/a/b/c/d"
	PrintSplitPath(p1, "a")
	PrintSplitPath(p2, "a")
	PrintSplitPath("/l", "a")
}
