package util

import (
	"github.com/czy21/cloud-disk-sync/exception"
	http2 "github.com/czy21/cloud-disk-sync/http"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
)

type HttpUtil struct {
	Request *resty.Request
}

func (h HttpUtil) Get(url string, v interface{}) error {

	res, err := h.Request.Get(url)
	exception.Check(err)
	var errMsg string
	if res.IsError() {
		errMsg = string(res.Body())
	}
	logParam := gin.LogFormatterParams{StatusCode: res.StatusCode(), Method: http.MethodGet, Path: url}
	log.Printf("| %3d | %13v |%-7s %#v\n%s",
		logParam.StatusCode,
		res.Time(),
		logParam.Method,
		logParam.Path,
		errMsg)
	err = http2.Client.JSONUnmarshal(res.Body(), v)
	exception.Check(err)
	return err
}
