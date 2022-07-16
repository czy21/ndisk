package http

import (
	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func Boot() {
	client = resty.New()
}

func GetClient() *resty.Client {
	return client
}
