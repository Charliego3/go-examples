package main

import (
	"database/sql/driver"
	"errors"
	"github.com/transerver/commons/utils"
)

// Bool type for db
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

	v0 := v[0]
	if v0 == 't' || v0 == 'T' {
		*b = true
	} else if v0 == 'f' || v0 == 'F' {
		*b = false
	} else {
		*b = v[0] >= 1
	}

	av := string(v)
	if utils.NotBlank(av) {
		if av == "true" || av == "TRUE" {
			*b = true
		} else if av == "false" || av == "FALSE" {
			*b = false
		}
		return nil
	}
	return nil
}
