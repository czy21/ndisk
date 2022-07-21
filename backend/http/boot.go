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
	// client.SetLogger(log.StandardLogger())
	// if viper.GetString("log.level") == "debug" {
	// 	client.SetDebug(true)
	// }
}

func GetClient() *resty.Client {
	return client
}
