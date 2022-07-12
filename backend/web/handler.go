package web

import (
	"fmt"
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/czy21/cloud-disk-sync/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case exception.MessageModel:
					c.JSON(http.StatusOK, model.ResponseModel{Error: err}.Build())
					break
				}
				eModel := exception.MessageModel{Code: fmt.Sprint(http.StatusInternalServerError), Message: fmt.Sprint(err)}
				c.JSON(http.StatusOK, model.ResponseModel{Error: eModel}.Build())
				panic(err)
			}
		}()
		c.Next()
	}
}

func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
