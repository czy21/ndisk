package model

import (
	"fmt"
	"strings"
	"time"
)

const StandardFormat = "2006-01-02 15:04:05"

type StandardTime time.Time

func (t *StandardTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*t).Format(StandardFormat) + `"`), nil
}

func (t *StandardTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}
	s, err := time.Parse(StandardFormat, value)
	if err != nil {
		return err
	}
	*t = StandardTime(s)
	return nil
}

// UnixTime unix timestamp
type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).UnixMilli())), nil
}
