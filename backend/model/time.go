package model

import (
	"fmt"
	"strings"
	"time"
)

const StandardFormat = "2006-01-02 15:04:05"

// StandardTime yyyy-MM-dd HH:mm:ss
type StandardTime time.Time

func (t StandardTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(StandardFormat)), nil
}

func (t *StandardTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	nt, err := time.ParseInLocation(StandardFormat, strings.Trim(string(data), `"`), time.Local)
	*t = StandardTime(nt)
	return err
}

// UnixTime unix timestamp
type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).UnixMilli())), nil
}
