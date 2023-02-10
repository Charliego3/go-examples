package logger

import (
	"log"
	"testing"
)

func TestCustomerLogrus(t *testing.T) {
	SetLevel(DebugLevel)
	// SetFormatter(&TextFormatter{
	//	FullTimestamp:   true,
	//	TimestampFormat: "2006-01-02 15:04:05.000000",
	//	PadLevelText:    true,
	// })
	SetFormatter(&Formatter{})

	log.SetFlags(0)
	log.SetOutput(StandardLogger().Writer())

	Info("Test 1")
	WithFields(Fields{
		"Application": nil,
		"number":      122,
	}).Error("A group of walrus emerges")

	WithFields(Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	WithField("Services", nil).Info("The group's number increased tremendously!")

	WithFields(Fields{
		"omg":    true,
		"number": 100,
	}).Debug("The ice breaks!")

	contextLogger := WithFields(Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

	log.Println("This is standard logger....")
}
