package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const method = "POST"
const contentType = "application/json"
const resource = "/api/logs"

// LogAnalyticsHook to send logs via the Telegram API.
type LogAnalyticsHook struct {
	CustomerID  string
	SharedKey   string
	apiEndpoint string
	LogType     string
	c           *http.Client
	async       bool
}

// Config defines a method for additional configuration when instantiating TelegramHook
type Config func(*LogAnalyticsHook)

// NewLogAnalyticsHook ...
func NewLogAnalyticsHook(CustomerID, SharedKey, LogType string, config ...Config) (*LogAnalyticsHook, error) {
	client := &http.Client{}
	return NewLogAnalyticsHookWithClient(CustomerID, SharedKey, LogType, client, config...)
}

// NewLogAnalyticsHookWithClient ...
func NewLogAnalyticsHookWithClient(CustomerID, SharedKey, LogType string, client *http.Client, config ...Config) (*LogAnalyticsHook, error) {
	apiEndpoint := fmt.Sprintf(
		"https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01",
		CustomerID,
	)

	h := LogAnalyticsHook{
		CustomerID:  CustomerID,
		c:           client,
		SharedKey:   SharedKey,
		apiEndpoint: apiEndpoint,
		async:       false,
	}

	for _, c := range config {
		c(&h)
	}

	// TODO: Validate token

	return &h, nil
}

// BuildSignature ...
func BuildSignature(customerID, sharedKey, date, contentLength, method, contentType, resource string) (string, error) {

	var keyBytes []byte
	var err error

	n := "\n"
	xHeaders := "x-ms-date:" + date
	signature := method + n + contentLength + n + contentType + n + xHeaders + n + resource
	bytesToHash := []byte(signature)

	if keyBytes, err = base64.StdEncoding.DecodeString(sharedKey); err != nil {
		fmt.Println("decode error:", err)
		return "", err
	}

	mac := hmac.New(sha256.New, keyBytes)
	mac.Write(bytesToHash)
	sum := mac.Sum(nil)
	sumEncoded := base64.StdEncoding.EncodeToString(sum)
	authorization := "SharedKey " + customerID + ":" + sumEncoded
	return authorization, nil
}

// PostLog ...
func (hook *LogAnalyticsHook) PostLog(messageJSONBytes []byte) error {

	date := time.Now().UTC().Format(time.RFC1123)
	length := strconv.Itoa(len(messageJSONBytes))
	signature, err := BuildSignature(hook.CustomerID, hook.SharedKey, date, length, method, contentType, resource)

	req, err := http.NewRequest(method, hook.apiEndpoint, bytes.NewReader(messageJSONBytes))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", signature)
	req.Header.Add("Log-Type", hook.LogType)
	req.Header.Add("x-ms-date", date)

	if res, err := hook.c.Do(req); err != nil {
		print("Failed with status code: " + string(res.StatusCode))
		return err
	}

	return nil
}

// Fire ...
func (hook *LogAnalyticsHook) Fire(entry *logrus.Entry) error {

	var log []byte
	var err error

	if log, err = json.MarshalIndent(entry.Data, "", "\t"); err != nil {
		print(err)
		return err
	}

	hook.PostLog(log)
	return nil
}

// Levels ...
func (hook *LogAnalyticsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
