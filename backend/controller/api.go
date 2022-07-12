package controller

import (
	"github.com/czy21/cloud-disk-sync/web"
	"github.com/czy21/cloud-disk-sync/webdav"
	"github.com/gin-gonic/gin"
)

func ApiEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: web.LogFormatter("API"),
	}))
	r.Use(gin.Recovery())
	r.Use(web.ErrorHandle())
	OptionController(r)
	webdav.WebDavController(r)
	return r
}
