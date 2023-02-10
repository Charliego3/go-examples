package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	SetLevel(DebugLevel)
	logger := NewLogger(WithPrefix("TestLogger"))
	logger.SetFormatter(&Formatter{})
	// logger.AddHook(&PrefixHook{})
	logger.WithField("field 1", "value 1").Error("withField error")
	logger.WithFields(Fields{
		"fields1": "value1",
		"fields2": "value2",
	}).Debug("withFields debug")
	logger.Info("This is 1")
	logger.Warnf("A group of walrus emerges %s", "args")

	logger = NewLogger()
	logger.Info("This is no Prefix")
	logger.Warnf("A group of walrus emerges %s", "args")
}
