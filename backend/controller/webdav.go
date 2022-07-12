package controller

import (
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/web"
	"github.com/gin-gonic/gin"
)

func WebDavTest(c *gin.Context) {
	input := model.OptionQuery{}
	err := c.Bind(&input)
	exception.Check(err)
	web.Context{Context: c}.
		OK(model.ResponseModel{Data: "hahaha"}.Build())

}

func WebDavController(r *gin.Engine) {

	v1 := r.Group("/dav")
	{
		v1.POST("/test", WebDavTest)
	}

}
