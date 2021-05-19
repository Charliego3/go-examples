package main

import (
	"database/sql/driver"
	"errors"
	errors2 "github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

type Bool bool

func (b Bool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

// Scan implements the sql.Scanner interface,
// and turns the bitfield incoming from MySQL into a BitBool
func (b *Bool) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion for Bool")
	}
	*b = v[0] == 1
	return nil
}

type Time time.Time

func (t Time) Value() (driver.Value, error) {
	return t, nil
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion for Time")
	}
	pt, err := time.Parse("2006-01-02 15:04:05", string(v))
	if err != nil {
		return errors2.Wrap(err, "时间格式化失败")
	}
	*t = Time(pt)
	return nil
}

type BigDecimal decimal.Decimal

func (t BigDecimal) String() string {
	return decimal.Decimal(t).String()
}

func (t BigDecimal) Value() (driver.Value, error) {
	return t, nil
}

func (t *BigDecimal) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	v, ok := src.([]byte)
	if !ok {
		return errors.New("bad []byte type assertion for Time")
	}
	dec, err := decimal.NewFromString(string(v))
	if err != nil {
		return errors2.Wrap(err, "can not parse for BigDecimal")
	}
	*t = BigDecimal(dec)
	return nil
}
