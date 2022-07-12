package webdav

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func Controller(r *gin.Engine) {
	v1 := r.Group("/dav")
	{
		v1.Any("/*path", ServeWebDAV)
		v1.Any("", ServeWebDAV)
		v1.Handle("PROPFIND", "/*path", ServeWebDAV)
		v1.Handle("PROPFIND", "", ServeWebDAV)
		v1.Handle("MKCOL", "/*path", ServeWebDAV)
		v1.Handle("LOCK", "/*path", ServeWebDAV)
		v1.Handle("UNLOCK", "/*path", ServeWebDAV)
		v1.Handle("PROPPATCH", "/*path", ServeWebDAV)
		v1.Handle("COPY", "/*path", ServeWebDAV)
		v1.Handle("MOVE", "/*path", ServeWebDAV)
	}
}
func ServeWebDAV(c *gin.Context) {
	handler := webdav.Handler{Prefix: "/dav", FileSystem: CloudFileSystem{}, LockSystem: webdav.NewMemLS()}
	handler.ServeHTTP(c.Writer, c.Request)
}
