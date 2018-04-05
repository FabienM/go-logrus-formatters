package formatters

import (
	"bytes"
	"encoding/json"
	"log/syslog"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGelfFormatter_Format(t *testing.T) {
	log := logrus.New()
	log.Formatter = NewGelf("testhost")
	buffer := new(bytes.Buffer)
	log.Out = buffer
	log.WithField("foo","bar").Info("great test message")

	var message map[string]interface{}
	err:= json.Unmarshal(buffer.Bytes(), &message)
	if err != nil {
		t.Error(err)
	}

	expectations:= map[string]interface{}{
		"_level_name": "INFORMATIONAL",
		"_foo": "bar",
		"host": "testhost",
		"level": float64(syslog.LOG_INFO),
		"short_message": "great test message",
		"version": GelfVersion,
	}
	for key, expected := range expectations {
		if message[key] != expected {
			t.Errorf("invalid log object: expected value for key '%s' was '%s', found '%s'", key, expected, message[key])
		}
	}
}
