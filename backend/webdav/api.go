package webdav

import (
	"context"
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/provider/local"
	"github.com/czy21/ndisk/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
	"strings"
)

var providerMetas []model.ProviderFolderMeta

func Controller(r *gin.Engine) {
	davPrefix := "/dav"
	serveFn := func(c *gin.Context) {
		providerMetas = repository.Provider{}.SelectListMeta()
		h := webdav.Handler{
			Prefix:     davPrefix,
			FileSystem: FileSystem{},
			LockSystem: webdav.NewMemLS(),
		}
		if c.Request.Method == "GET" {
			_, fs := getProvider(strings.TrimPrefix(c.Request.URL.Path, davPrefix), "")
			switch fs.(type) {
			case local.FileSystem:
				break
			default:
				DownloadFile(c.Writer, c.Request, strings.TrimPrefix(c.Request.URL.Path, davPrefix))
				return
			}
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "Request", c.Request))
		h.ServeHTTP(Writer{c.Writer}, c.Request)
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
