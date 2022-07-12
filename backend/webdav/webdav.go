package webdav

import (
	"github.com/czy21/cloud-disk-sync/constant"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func Controller(r *gin.Engine) {
	v1 := r.Group("/dav")
	{
		v1.Any("/*path", ServeWebDAV)
		for _, t := range constant.WebDavMethods {
			v1.Handle(t, "/*path", ServeWebDAV)
		}
	}
}
func ServeWebDAV(c *gin.Context) {
	handler := webdav.Handler{Prefix: "/dav", FileSystem: webdav.Dir("./data/"), LockSystem: webdav.NewMemLS()}
	handler.ServeHTTP(c.Writer, c.Request)
}
