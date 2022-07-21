package webdav

import (
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
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
			DownloadFile(c.Writer, c.Request, strings.TrimPrefix(c.Request.URL.Path, davPrefix))
		} else {
			h.ServeHTTP(c.Writer, c.Request)
		}
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
