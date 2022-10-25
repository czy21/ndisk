package webdav

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var p1 = ""
var p2 = "/"
var p3 = "/a"
var p4 = "/a/b/c"
var p5 = "/a/b/c/"

func splitResult(p string) {
	d, f := filepath.Split(p)
	d1, f1 := path.Split(p)
	fmt.Println("============================")
	fmt.Println(fmt.Sprintf("filepath.Split path:%s dir:%s file:%s", p, d, f))
	fmt.Println(fmt.Sprintf("    path.Split path:%s dir:%s file:%s", p, d1, f1))
	fmt.Println("============================")
}

func TestSplit(t *testing.T) {
	splitResult(p1)
	splitResult(p2)
	splitResult(p3)
	splitResult(p4)
	splitResult(p5)
}

func dirResult(p string) {
	d := filepath.Dir(p)
	b := filepath.Base(p)
	fmt.Println("============================")
	fmt.Println(fmt.Sprintf("path.Dir  path:%s ret:%s", p, d))
	fmt.Println(fmt.Sprintf("path.Base path:%s ret:%s", p, b))
	fmt.Println("============================")
}

func TestDirBase(t *testing.T) {
	dirResult(p1)
	dirResult(p2)
	dirResult(p3)
	dirResult(p4)
	dirResult(p5)
	a := "/a123"
	b := path.Join(strings.SplitAfter(a, "/")[0:2]...)
	c := strings.Split("a/b/c/", "/")
	d := c[len(c)-1:]
	println(d)
	println(b)
}
