package model

import (
	"time"
)

type ResponseModel struct {
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp UnixTime    `json:"timestamp"`
}

func (r ResponseModel) Build() ResponseModel {
	r.Timestamp = UnixTime(time.Now())
	return r
}
