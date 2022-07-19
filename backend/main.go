package main

import (
	"fmt"
	"github.com/czy21/ndisk/bootstrap"
	"os"
	"time"
)

func main() {
	tz, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		fmt.Println(err)
	}
	time.Local = tz
	bootstrap.Boot()
}
