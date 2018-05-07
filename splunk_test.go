package splunk

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBuildSignature(t *testing.T) {

	customerID := os.Getenv("LA_CUST_ID")
	sharedKey := os.Getenv("LA_SHARED_KEY")
	TimeStampField := "Tue, 01 May 2018 23:10:18 UTC"
	ContentType := "application/json"
	resource := "/api/logs"
	length := "60"
	method := "POST"

	// Remember to update hash if LA_CUST_ID or LA_SHARED_KEY change
	expected := "SharedKey " + customerID + ":M/YcOOY2e792N7ebR9fW5ZeJ+oXpVJN/7JZi79+2Kq8="

	resp, _ := BuildSignature(customerID, sharedKey, TimeStampField, length, method, ContentType, resource)

	assert.Equal(t, resp, expected, "Two signatures should be the same")
}

func TestNewHook(t *testing.T) {

	customerID := os.Getenv("LA_CUST_ID")
	sharedKey := os.Getenv("LA_SHARED_KEY")

	h, err := NewLogAnalyticsHook(customerID, sharedKey, "IntegrationTest")

	if err == nil {
		t.Errorf("No error on invalid Telegram API token.")
	}

	log.AddHook(h)

	log.WithFields(log.Fields{
		"animal":   "walrus",
		"number":   1,
		"size":     10,
		"log_type": "ItWorks",
	}).Errorf("A walrus appears")
}
