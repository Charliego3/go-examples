package times

import (
	"strconv"
	"time"
)

func Parse(millis int64) time.Time {
	return time.Unix(0, millis*time.Millisecond.Nanoseconds())
}

func Parse2S(millis int64, format ...string) string {
	formater := "2006-06-01 15:04:05.000"
	if len(format) > 0 {
		formater = format[0]
	}
	return Parse(millis).Format(formater)
}

func Parse2I2S(millis string, format ...string) string {
	i, err := strconv.ParseInt(millis, 10, 0)
	if err != nil {
		return millis
	}
	return Parse2S(i, format...)
}
