package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/strutil"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"runtime"
	"time"
)

func init() {
	SetupLogger(StandardLogger())
}

type Config struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool

	// //// Text RusFormatter
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool
}

type Logger struct {
	*logrus.Logger
	Prefix string
}

type Option interface{ apply(logger *Logger) }
type prefix string
type formatter struct{ RusFormatter }

func (f formatter) apply(logger *Logger) { logger.Formatter = f.RusFormatter }
func (p prefix) apply(logger *Logger)    { logger.Prefix = string(p) }

func WithPrefix(format string, v ...interface{}) Option { return prefix(fmt.Sprintf(format, v...)) }
func WithFormatter(f RusFormatter) Option               { return formatter{f} }

func NewLogger(opts ...Option) *Logger {
	logger := &Logger{Logger: StandardLogger()}
	for _, opt := range opts {
		opt.apply(logger)
	}
	return logger
}

func SetupLogger(logger *logrus.Logger) {
	// if gin.IsDebugging() {
	logger.SetLevel(DebugLevel)
	logger.SetFormatter(&Formatter{})
	// } else {
	// 	logger.SetFormatter(&JSONFormatter{
	// 		TimestampFormat: "2006-01-02 15:04:05.000000",
	// 	})
	// }

	gin.DefaultErrorWriter = logger.Writer()
	gin.DefaultWriter = logger.Writer()
	log.SetFlags(0)
	log.SetOutput(logger.Writer())
}

func (l *Logger) Log(level Level, args ...interface{}) {
	l.withField().Log(level, args...)
}

func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	l.withField().Logf(level, format, args...)
}

func (l *Logger) withField() *Entry {
	return l.WithField("", nil)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (l *Logger) WithField(key string, value interface{}) *Entry {
	if len(key) == 0 && value == nil {
		if strutil.IsNotBlank(l.Prefix) {
			return l.Logger.WithFields(Fields{l.Prefix: nil, key: value})
		}
		return logrus.NewEntry(l.Logger)
	}

	if value != nil {

	}

	if len(l.Prefix) == 0 {
		if len(key) == 0 && value == nil {
			return logrus.NewEntry(l.Logger)
		}
		return l.Logger.WithField(key, value)
	} else if value == nil {
		prefix := l.Prefix
		if len(key) != 0 {
			prefix = key
		}
		return l.Logger.WithField(prefix, nil)
	} else {
		return l.Logger.WithFields(Fields{l.Prefix: nil, key: value})
	}
}

// WithFields Adds a struct of fields to the log entry. All it does is call `WithField` for each `Field`.
func (l *Logger) WithFields(fields Fields) *Entry {
	if len(l.Prefix) != 0 {
		fields[l.Prefix] = nil
	}
	return l.Logger.WithFields(fields)
}

// WithPrefix Adds a struct of fields to the log entry. All it does is call `WithField` for each `Field`.
func (l *Logger) WithPrefix(prefix string) *Entry {
	return l.WithField(prefix, nil)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Logf(TraceLevel, format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logf(DebugLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logf(InfoLevel, format, args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logf(WarnLevel, format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logf(FatalLevel, format, args...)
	l.Exit(1)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logf(PanicLevel, format, args...)
}

func (l *Logger) Trace(args ...interface{}) {
	l.Log(TraceLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Log(DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Log(InfoLevel, args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Log(WarnLevel, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Log(ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Log(FatalLevel, args...)
	l.Exit(1)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Log(PanicLevel, args...)
}

func (l *Logger) Logln(level Level, args ...interface{}) {
	l.withField().Logln(level, args...)
}

func (l *Logger) Traceln(args ...interface{}) {
	l.Logln(TraceLevel, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Logln(DebugLevel, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Logln(InfoLevel, args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.Infoln(args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Logln(WarnLevel, args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.Warnln(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Logln(ErrorLevel, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Logln(FatalLevel, args...)
	l.Exit(1)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.Logln(PanicLevel, args...)
}

var (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
	AllLevels  = logrus.AllLevels
)

type (
	Level         = logrus.Level
	RusFormatter  = logrus.Formatter
	Hook          = logrus.Hook
	Entry         = logrus.Entry
	Fields        = logrus.Fields
	FieldMap      = logrus.FieldMap
	TextFormatter = logrus.TextFormatter
	JSONFormatter = logrus.JSONFormatter
)

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter RusFormatter) {
	logrus.SetFormatter(formatter)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	logrus.SetReportCaller(include)
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	return logrus.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level Level) bool {
	return logrus.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return logrus.WithField(logrus.ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *Entry {
	return logrus.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return logrus.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return logrus.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *Entry {
	return logrus.WithTime(t)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logrus.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	logrus.Traceln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logrus.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logrus.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}
