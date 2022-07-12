package webdav

import (
	"fmt"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/web"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func WebDavTest(c *gin.Context) {
	input := model.OptionQuery{}
	err := c.Bind(&input)
	exception.Check(err)
	web.Context{Context: c}.
		OK(model.ResponseModel{Data: "hahaha"}.Build())

}

func WebDavController(r *gin.Engine) {
	r.Any("/dav/*proxyPath", func(c *gin.Context) {
		fmt.Println(c.Param("proxyPath"))
		handler := webdav.Handler{Prefix: "/dav", FileSystem: CloudFileSystem{}, LockSystem: webdav.NewMemLS()}
		handler.ServeHTTP(c.Writer, c.Request)
	})
}
