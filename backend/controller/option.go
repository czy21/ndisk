package controller

import (
	"context"
	"github.com/czy21/cloud-disk-sync/cache"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/czy21/cloud-disk-sync/service"
	"github.com/czy21/cloud-disk-sync/web"
	"github.com/gin-gonic/gin"
)

func OptionList(c *gin.Context) {
	input := model.OptionQuery{}
	err := c.Bind(&input)
	exception.Check(err)
	web.Context{Context: c}.
		OK(model.ResponseModel{Data: service.Option{}.FindByKeys(input)}.Build())

}

func CachePut(c *gin.Context) {
	input := make(map[string]interface{})
	err := c.Bind(&input)
	exception.Check(err)
	cache.Client.SetObj(context.Background(), c.Query("key"), input)
}

func OptionController(r *gin.Engine) {

	v1 := r.Group("/option")
	{
		v1.POST("/query", OptionList)
		v1.POST("/cache/put", CachePut)
	}

}
