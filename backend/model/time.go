package model

import (
	"fmt"
	"strings"
	"time"
)

const StandardFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*t).Format(StandardFormat) + `"`), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}
	s, err := time.ParseInLocation(StandardFormat, value, time.UTC)
	if err != nil {
		return err
	}
	*t = LocalTime(s)
	return nil
}

// UnixTime unix timestamp
type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).UnixMilli())), nil
}
