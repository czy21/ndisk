package webdav

import (
	"github.com/czy21/cloud-disk-sync/constant"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

var providerMetas []model.ProviderMeta

func Controller(r *gin.Engine) {
	providerMetas = repository.Provider{}.SelectList()
	serveFn := func(c *gin.Context) {
		h := webdav.Handler{
			Prefix:     "/dav",
			FileSystem: FileSystem{},
			LockSystem: webdav.NewMemLS(),
		}
		h.ServeHTTP(c.Writer, c.Request)
	}
	r1 := r.Group("/dav")
	{
		r1.Any("/*path", serveFn)
		for _, t := range constant.WebDavMethods {
			r1.Handle(t, "/*path", serveFn)
		}
	}
}
