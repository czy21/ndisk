package util

import (
	"github.com/czy21/cloud-disk-sync/exception"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
)

type HttpUtil struct {
	*resty.Client
}

func (HttpUtil) NewClient() HttpUtil {
	c := resty.New()
	return HttpUtil{Client: c}
}

func (h HttpUtil) Get(url string, v interface{}) {
	res, err := h.R().Get(url)
	exception.Check(err)
	var errMsg string
	if res.IsError() {
		errMsg = string(res.Body())
	}
	logParam := gin.LogFormatterParams{StatusCode: res.StatusCode(), Method: http.MethodGet, Path: url, ClientIP: h.BaseURL}
	log.Printf("| %3d | %13v | %15s |%-7s %#v\n%s",
		logParam.StatusCode,
		res.Time(),
		logParam.ClientIP,
		logParam.Method,
		logParam.Path,
		errMsg)
	err = h.JSONUnmarshal(res.Body(), v)
	exception.Check(err)
}
