package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func Boot() {
	gin.DefaultWriter = log.Writer()
	if viper.GetString("log.file") != "" {
		gin.DisableConsoleColor()
	}
	var g errgroup.Group
	g.Go(func() error {
		address := fmt.Sprintf(":%s", viper.GetString("server.port"))
		log.Printf("Listening and api HTTP on %s\n", address)
		err := http.ListenAndServe(address, ApiEngine())
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		address := fmt.Sprintf(":%s", viper.GetString("web.port"))
		log.Printf("Listening and web HTTP on %s\n", address)
		err := http.ListenAndServe(address, WebEngine())
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
