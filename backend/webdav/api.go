package webdav

import (
	"context"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
	"net/http"
	"strings"
)

var providerMetas []model.ProviderFolderMeta

func getDavLogger() func(request *http.Request, err error) {
	return func(request *http.Request, err error) {
		if err != nil {
			log.Debugf("%s %s", request.RequestURI, err)
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
		ctx := context.WithValue(request.Context(), "contentLength", request.ContentLength|0)
		request = request.WithContext(ctx)
		HandleHttp(strings.TrimPrefix(c.Request.URL.Path, davPrefix), &writer, request)
		h.ServeHTTP(writer, request)
	}
	r1 := r.Group(davPrefix)
	{
		r1.Any("/*path", serveFn)
		r1.Any("", serveFn)
		for _, t := range constant.WebDavMethods {
			r1.Handle(t, "/*path", serveFn)
		}
	}
}
