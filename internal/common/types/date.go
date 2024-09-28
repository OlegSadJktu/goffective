package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}


func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	dates := strings.Split(s, ".")

	day, err := strconv.Atoi(dates[0])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(dates[1])
	if err != nil {
		return err
	}
	year, err := strconv.Atoi(dates[2])
	if err != nil {
		return err
	}
  asd := time.UTC
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, asd)
	ct.Time = t
	return
}

func (ct *CustomTime) IsSet() bool {
	return !ct.IsZero()
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%v.%v.%v", ct.Day(), ct.Month(), ct.Year())), nil
}
