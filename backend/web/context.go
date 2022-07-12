package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context struct {
	*gin.Context
}

func (c Context) OK(obj any) {
	c.JSON(http.StatusOK, obj)
}
