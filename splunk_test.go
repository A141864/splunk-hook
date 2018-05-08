package splunk

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGetSessionKey(t *testing.T) {

	token := os.Getenv("token")

	defaultValues := log.Fields{
		"CustomField1": "Some Value",
		"CustomField2": "Another Value",
	}
	log.AddHook(NewHook(
		"https://127.0.0.1:8088/services/collector",
		token,
		"SourceApp",
		"json",
		"default",
		defaultValues,
	))

	log.WithFields(log.Fields{
		"CustomField": "it works!",
	}).Error("This is the error message")
}
