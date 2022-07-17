package util

import (
	"github.com/czy21/ndisk/exception"
	http2 "github.com/czy21/ndisk/http"
	"github.com/go-resty/resty/v2"
)

type HttpUtil struct {
	Request *resty.Request
}

func (h HttpUtil) Get(url string, v interface{}) error {
	res, err := h.Request.Get(url)
	exception.Check(err)
	//var errMsg string
	//if res.IsError() {
	//	errMsg = string(res.Body())
	//}
	//logParam := gin.LogFormatterParams{
	//	TimeStamp:    res.Request.Time,
	//	StatusCode:   res.StatusCode(),
	//	Method:       http.MethodGet,
	//	Path:         url,
	//	ErrorMessage: errMsg,
	//	Latency:      res.Time(),
	//}
	//log.Debugf("[RPC] %s", web.LogFormatter()(logParam))
	err = http2.GetClient().JSONUnmarshal(res.Body(), v)
	exception.Check(err)
	return err
}
