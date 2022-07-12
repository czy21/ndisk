package model

import (
	"fmt"
	"time"
)

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).UnixMilli())), nil
}

type ResponseModel struct {
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp UnixTime    `json:"timestamp"`
}

func (r ResponseModel) Build() ResponseModel {
	r.Timestamp = UnixTime(time.Now())
	return r
}
