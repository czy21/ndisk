package webdav

import (
	"github.com/czy21/ndisk/constant"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

var providerMetas []model.ProviderFolderMeta

func Controller(r *gin.Engine) {
	serveFn := func(c *gin.Context) {
		providerMetas = repository.Provider{}.SelectListMeta()
		davPrefix := "/dav"
		h := webdav.Handler{
			Prefix:     davPrefix,
			FileSystem: FileSystem{},
			LockSystem: webdav.NewMemLS(),
		}
		h.ServeHTTP(c.Writer, c.Request)
	}
	r1 := r.Group("/dav")
	{
		r1.Any("/*path", serveFn)
		r1.Any("", serveFn)
		for _, t := range constant.WebDavMethods {
			r1.Handle(t, "/*path", serveFn)
		}
	}
}
