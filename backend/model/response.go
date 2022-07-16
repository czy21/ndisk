package model

import (
	"fmt"
	"strings"
	"time"
)

type UnixTime time.Time
type StandardTime time.Time // yyyy-MM-dd HH:mm:ss

func (t *StandardTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	nt, err := time.ParseInLocation("2006-01-02 15:04:05", strings.Trim(string(data), `"`), time.Local)
	*t = StandardTime(nt)
	return err
}
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
