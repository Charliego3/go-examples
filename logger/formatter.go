package logger

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/strutil"
	"io"
	"strings"
	"sync"
)

type Formatter struct {
	Config         Config
	colors         map[Level]*color.Color
	colorOnce      sync.Once
	underLineColor *color.Color
}

func (f *Formatter) init() {
	f.colorOnce.Do(func() {
		f.Config.TimestampFormat = "2006-01-02 15:04:05.000000"
		if !gin.IsDebugging() {
			return
		}

		f.underLineColor = color.New(color.Underline)
		f.colors = make(map[Level]*color.Color)

		for _, l := range AllLevels {
			f.initLevelColor(l)
		}
	})
}

func (f *Formatter) Format(entry *Entry) ([]byte, error) {
	f.init()
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	cf := f.colors[entry.Level]
	f.write(cf, b, "[%s]", f.upperLevel(entry.Level))
	f.write(cf, b, " [%s]", entry.Time.Format(f.Config.TimestampFormat))

	var prefixLen int
	for k, v := range entry.Data {
		prefix := v == nil
		if !prefix {
			if s, ok := v.(string); ok && len(s) == 0 {
				prefix = true
			}
		}

		delete(entry.Data, k)

		if prefix && strutil.IsNotBlank(k) {
			var pf string
			if f.underLineColor == nil {
				pf = fmt.Sprintf("[%s]", k)
			} else {
				pf = f.underLineColor.Sprintf("[%s]", k)
			}
			f.write(cf, b, " %s:", pf)
			prefixLen = len(k) + 4
			break
		}
	}

	if len(entry.Message)+prefixLen > 50 || len(entry.Data) == 0 {
		f.write(cf, b, " %s", entry.Message)
	} else {
		f.write(cf, b, " %-*s", 50-prefixLen, entry.Message)
	}

	for k, v := range entry.Data {
		if strutil.IsBlank(k) && v == nil {
			continue
		}

		f.write(cf, b, " %s", k)
		f.write(cf, b, "=")
		f.appendValue(b, v)
	}

	return append(b.Bytes(), '\n'), nil
}

func (f *Formatter) write(cf *color.Color, w io.Writer, format string, a ...any) {
	if cf == nil {
		_, _ = fmt.Fprintf(w, format, a...)
	} else {
		_, _ = cf.Fprintf(w, format, a...)
	}
}

func (f *Formatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *Formatter) needsQuoting(text string) bool {
	if f.Config.ForceQuote {
		return true
	}
	if f.Config.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.Config.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *Formatter) upperLevel(level Level) string {
	switch level {
	case DebugLevel:
		return "DBUG"
	default:
		return strings.ToUpper(level.String()[:4])
	}
}

func (f *Formatter) initLevelColor(level Level) {
	var lc color.Attribute
	switch level {
	case DebugLevel, TraceLevel:
		lc = color.FgWhite
	case WarnLevel:
		lc = color.FgYellow
	case ErrorLevel, FatalLevel, PanicLevel:
		lc = color.FgRed
	default:
		lc = color.FgCyan
	}
	c := color.New(lc, color.Bold)
	f.colors[level] = c
}
