package controller

import (
	"fmt"
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

func CachePut(c *gin.Context) {
	t1 := time.Now()
	fmt.Println(fmt.Sprintf("t1: %s", t1))
	t2, _ := time.Parse(model.StandardFormat, "2022-07-19 12:00:00")
	fmt.Println(fmt.Sprintf("t2: %s", t2.Add(-8*time.Hour)))
	fmt.Println(time.Now())
}

func OptionController(r *gin.Engine) {

	v1 := r.Group("/api/option")
	{
		v1.POST("/query", OptionList)
		v1.POST("/cache/put", CachePut)
	}

}
