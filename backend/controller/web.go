package controller

import (
	"fmt"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/exception"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func backendReverseProxy(c *gin.Context, subPath string) {
	path := fmt.Sprintf("http://127.0.0.1:%s", viper.GetString("server.port"))
	remote, err := url.Parse(path)
	exception.Check(err)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = subPath + c.Param("path")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func WebEngine() *gin.Engine {
	r := gin.New()
	//r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
	//	Formatter: web.LogFormatter("WEB"),
	//}))
	r.Use(gin.Recovery())
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		indexFile := fmt.Sprintf("%s/index.html", viper.GetString("web.dist"))
		staticFile := fmt.Sprintf("%s/static/", viper.GetString("web.dist"))
		r.NoRoute(func(c *gin.Context) {
			c.File(indexFile)
		})
		r.StaticFile("/", indexFile)
		r.Static("/static/", staticFile)
	}
	r.Any("/api/*path", func(c *gin.Context) {
		backendReverseProxy(c, "")
	})
	r.Any("/dav/*path", func(c *gin.Context) {
		backendReverseProxy(c, "/dav")
	})
	for _, t := range constant.WebDavMethods {
		r.Handle(t, "/dav/*path", func(c *gin.Context) {
			backendReverseProxy(c, "/dav")
		})
	}
	return r
}
