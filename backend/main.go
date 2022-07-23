package main

import (
	"github.com/czy21/ndisk/bootstrap"
	"time"
)

func main() {
	time.Local = time.UTC
	bootstrap.Boot()
}
