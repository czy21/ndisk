package controller

import (
	"github.com/czy21/ndisk/exception"
	"github.com/czy21/ndisk/model"
	"github.com/czy21/ndisk/service"
	"github.com/czy21/ndisk/web"
	"github.com/gin-gonic/gin"
	"time"
)

func OptionList(c *gin.Context) {
	input := model.OptionQuery{}
	err := c.Bind(&input)
	exception.Check(err)
	web.Context{Context: c}.
		OK(model.ResponseModel{Data: service.Option{}.FindByKeys(input)}.Build())

}

type TestVO struct {
	T model.UnixTime `json:"t"`
}

func CachePut(c *gin.Context) {
	var input TestVO
	_ = c.Bind(&input)
	web.Context{Context: c}.OK(model.ResponseModel{Data: map[string]interface{}{
		"t": time.Now(),
		"i": input.T,
	}}.Build())
}

func OptionController(r *gin.Engine) {

	v1 := r.Group("/api/option")
	{
		v1.POST("/query", OptionList)
		v1.POST("/cache/put", CachePut)
	}

}
