package http

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func Boot() {
	client = resty.New()
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal
}

func GetClient() *resty.Client {
	return client
}
