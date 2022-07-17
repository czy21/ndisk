package util

import (
	"github.com/czy21/ndisk/exception"
	"github.com/go-resty/resty/v2"
)

type HttpUtil struct {
	Request *resty.Request
}

func (h HttpUtil) Get(url string, v interface{}) error {
	h.Request.SetResult(v)
	_, err := h.Request.Get(url)
	exception.Check(err)
	return err
}
