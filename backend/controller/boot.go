package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func Boot() {
	gin.SetMode(gin.ReleaseMode)
	logLevel := viper.GetString("log.level")
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	address := fmt.Sprintf(":%s", viper.GetString("server.port"))
	log.Infof("TZ: %s Date: %s", time.Local, time.Now())
	log.Infof("Listening and serving HTTP on %s", address)
	err := http.ListenAndServe(address, ApiEngine())
	if err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
}
