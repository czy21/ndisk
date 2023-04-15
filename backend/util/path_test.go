package util

import (
	"fmt"
	"testing"
)

func PrintSplitPath(p string, pPrefix string) {
	dir, base, parents, rel, isRoot := SplitPath(p, pPrefix)
	fmt.Println(fmt.Sprintf("%s dir:%s base:%s parents:%s rel:%s isRoot:%t", p, dir, base, parents, rel, isRoot))
}

func TestSplitPath(t *testing.T) {
	p1 := "a/b/c/d"
	p2 := "/a/b/c/d"
	PrintSplitPath(p1, "a")
	PrintSplitPath(p2, "a")
	PrintSplitPath("/l", "a")
}
