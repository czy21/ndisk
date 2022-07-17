package controller

import (
	"github.com/czy21/ndisk/web"
	"github.com/czy21/ndisk/webdav"
	"github.com/gin-gonic/gin"
)

func ApiEngine() *gin.Engine {
	r := gin.New()
	r.Use(web.LogHandler())
	r.Use(gin.Recovery())
	r.Use(web.ErrorHandler())
	OptionController(r)
	webdav.Controller(r)
	return r
}
