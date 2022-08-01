package webdav

import (
	"context"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/webdav"
	"net/http"
	"strings"
)

var providerMetas []model.ProviderFolderMeta

func getDavLogger() func(request *http.Request, err error) {
	return func(request *http.Request, err error) {
		if err != nil {
			log.Errorf("%s %s", request.RequestURI, err)
		}
	}
}

func Controller(r *gin.Engine) {
	davPrefix := "/dav"
	serveFn := func(c *gin.Context) {
		providerMetas = repository.Provider{}.SelectListMeta()
		h := webdav.Handler{
			Prefix:     davPrefix,
			FileSystem: FileSystem{},
			LockSystem: webdav.NewMemLS(),
		}
		h.Logger = getDavLogger()
		var writer http.ResponseWriter = c.Writer
		var request = c.Request
		extra := map[string]interface{}{
			constant.HttpExtraFileSize: request.ContentLength | 0,
			constant.HttpExtraMethod:   request.Method,
		}
		ctx := context.WithValue(request.Context(), constant.HttpExtra, extra)
		request = request.WithContext(ctx)
		HandleHttp(strings.TrimPrefix(c.Request.URL.Path, davPrefix), &writer, request)
		h.ServeHTTP(writer, request)
	}
	r1 := r.Group(davPrefix)
	{
		r1.Use(BasicAuth)
		r1.Any("/*path", serveFn)
		r1.Any("", serveFn)
		for _, t := range constant.WebDavMethods {
			r1.Handle(t, "/*path", serveFn)
		}
	}
}

func BasicAuth(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Writer.Header()["WWW-Authenticate"] = []string{`Basic realm="ndisk"`}
		c.Status(http.StatusUnauthorized)
		c.Abort()
		return
	}
	if viper.GetString("dav.username") == username && viper.GetString("dav.password") == password {
		c.Next()
		return
	}
	c.Status(http.StatusUnauthorized)
	c.Abort()
}
