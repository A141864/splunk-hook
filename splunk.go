package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Event ...
type Event struct {
	Time       int64       `json:"time" binding:"required"`       // epoch time in seconds
	Host       string      `json:"host" binding:"required"`       // hostname
	Source     string      `json:"source" binding:"required"`     // app name
	SourceType string      `json:"sourcetype" binding:"required"` // Splunk bucket to group logs in
	Index      string      `json:"index" binding:"required"`      // idk what it does..
	Event      interface{} `json:"event" binding:"required"`      // throw any useful key/val pairs here
}

// Hook ...
type Hook struct {
	HTTPClient    *http.Client // HTTP client used to communicate with the API
	URL           string
	Hostname      string
	Token         string
	Source        string //Default source
	SourceType    string //Default source type
	Index         string //Default index
	DefaultValues logrus.Fields
}

// Config defines a method for additional configuration when instantiating TelegramHook
type Config func(*Hook)

// Fire ...
func (s *Hook) Fire(entry *logrus.Entry) error {
	var log []byte
	var err error

	combinedFields := entry.Data

	switch entry.Level {
	case logrus.DebugLevel:
		combinedFields["logLevel"] = "debug"
	case logrus.InfoLevel:
		combinedFields["logLevel"] = "info"
	case logrus.ErrorLevel:
		combinedFields["logLevel"] = "error"
	case logrus.FatalLevel:
		combinedFields["logLevel"] = "fatal"
	case logrus.PanicLevel:
		combinedFields["logLevel"] = "panic"
	default:
		combinedFields["logLevel"] = "info"
	}

	combinedFields["message"] = entry.Message

	for k, v := range s.DefaultValues {
		combinedFields[k] = v
	}

	if log, err = json.MarshalIndent(combinedFields, "", "\t"); err != nil {
		print(err)
		return err
	}

	// print(string(log))
	s.Log(log)
	return nil
}

// Levels ...
func (s *Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// NewHook ...
func NewHook(URL string, Token string, Source string, SourceType string, Index string, defaultFields logrus.Fields) *Hook {

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Timeout: time.Second * 20, Transport: tr}

	hostname, _ := os.Hostname()
	c := Hook{
		HTTPClient:    httpClient,
		URL:           URL,
		Hostname:      hostname,
		Token:         Token,
		Source:        Source,
		SourceType:    SourceType,
		Index:         Index,
		DefaultValues: defaultFields,
	}

	return &c
}

// NewEvent ...
func (s *Hook) NewEvent(fields string, source string, sourcetype string, index string) *Event {

	e := &Event{
		Time:       time.Now().Unix(),
		Host:       s.Hostname,
		Source:     source,
		SourceType: sourcetype,
		Index:      index,
		Event:      fields,
	}
	return e
}

// Log ...
func (s *Hook) Log(fields interface{}) error {

	f := string(fields.([]byte))
	log := s.NewEvent(f, s.Source, s.SourceType, s.Index)
	return s.LogEvent(log)
}

// LogEvent ...
func (s *Hook) LogEvent(e *Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return s.doRequest(bytes.NewBuffer(b))
}

func (s *Hook) doRequest(b *bytes.Buffer) error {
	// make new request
	url := s.URL
	req, err := http.NewRequest("POST", url, b)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+s.Token)

	// receive response
	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	// If statusCode is not good, return error string
	switch res.StatusCode {
	case 200:
		return nil
	default:
		// Turn response into string and return it
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = errors.New(responseBody)

	}
	return err
}
