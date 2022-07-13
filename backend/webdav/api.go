package webdav

import (
	"github.com/czy21/cloud-disk-sync/constant"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func Controller(r *gin.Engine) {
	r1 := r.Group("/dav")
	{
		r1.Any("/*path", ServeHttp)
		for _, t := range constant.WebDavMethods {
			r1.Handle(t, "/*path", ServeHttp)
		}
	}
}
func ServeHttp(c *gin.Context) {
	h := webdav.Handler{
		Prefix:     "/dav",
		FileSystem: FileSystem{},
		LockSystem: webdav.NewMemLS(),
	}
	h.ServeHTTP(c.Writer, c.Request)
}
