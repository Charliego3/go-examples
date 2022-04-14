package main

import (
	"database/sql/driver"
	"errors"
	errors2 "github.com/pkg/errors"
	"time"
)

// DateTime type for db datetime
type DateTime struct {
	T *time.Time
}

func (t DateTime) Time() time.Time {
	return *t.T
}

func (t DateTime) Value() (driver.Value, error) {
	return t, nil
}

func (t *DateTime) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.([]byte)
	if ok && len(v) > 0 {
		pt, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return errors2.Wrap(err, "时间格式化失败")
		}
		*t = DateTime{T: &pt}
		return nil
	}

	v1, ok1 := src.(time.Time)
	if ok1 {
		*t = DateTime{&v1}
		return nil
	}
	return errors.New("bad []byte or time.Time type assertion for DateTime")
}

func (t *DateTime) String() string {
	if t == nil {
		return ""
	}
	return t.T.Format("2006-01-02 15:04:05.000000")
}
