package webdav

import (
	"context"
	"github.com/czy21/cloud-disk-sync/constant"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

var providerMetas []model.ProviderMeta

func Controller(r *gin.Engine) {
	serveFn := func(c *gin.Context) {
		providerMetas = repository.Provider{}.SelectList()
		h := webdav.Handler{
			Prefix:     "/dav",
			FileSystem: FileSystem{},
			LockSystem: webdav.NewMemLS(),
		}
		ctx := context.WithValue(c.Request.Context(), "env", make(map[string]interface{}))
		c.Request = c.Request.WithContext(ctx)
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
