package formatters

import (
	"encoding/json"
	"log/syslog"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// GelfVersion is the supported gelf version
	GelfVersion = "1.1"
)

var (
	levelMap      map[logrus.Level]syslog.Priority
	syslogNameMap map[syslog.Priority]string

	protectedFields map[string]bool

	// DefaultLevel is the default syslog level to use if the logrus level does not map to a syslog level
	DefaultLevel = syslog.LOG_INFO
)

func init() {
	levelMap = map[logrus.Level]syslog.Priority{
		logrus.PanicLevel: syslog.LOG_EMERG,
		logrus.FatalLevel: syslog.LOG_CRIT,
		logrus.ErrorLevel: syslog.LOG_ERR,
		logrus.WarnLevel:  syslog.LOG_WARNING,
		logrus.InfoLevel:  syslog.LOG_INFO,
		logrus.DebugLevel: syslog.LOG_DEBUG,
	}
	syslogNameMap = map[syslog.Priority]string{
		syslog.LOG_EMERG:   "EMERGENCY",
		syslog.LOG_ALERT:   "ALERT",
		syslog.LOG_CRIT:    "CRITICAL",
		syslog.LOG_ERR:     "ERROR",
		syslog.LOG_WARNING: "WARNING",
		syslog.LOG_NOTICE:  "NOTICE",
		syslog.LOG_INFO:    "INFORMATIONAL",
		syslog.LOG_DEBUG:   "DEBUGGING",
	}
	protectedFields = map[string]bool{
		"version":       true,
		"host":          true,
		"short_message": true,
		"full_message":  true,
		"timestamp":     true,
		"level":         true,
	}
}

type gelfFormatter struct {
	hostname string
}

// NewGelf returns a new logrus / gelf-compliant formatter
func NewGelf(hostname string) gelfFormatter {
	return gelfFormatter{hostname: hostname}
}

// Format implements logrus formatter
func (f gelfFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := toSyslogLevel(entry.Level)
	gelfEntry := map[string]interface{}{
		"version":       GelfVersion,
		"short_message": entry.Message,
		"level":         level,
		"timestamp":     toTimestamp(entry.Time),
		"host":          f.hostname,
		"_level_name":   syslogNameMap[level],
	}
	if _, file, line, ok := runtime.Caller(5); ok {
		gelfEntry["_file"] = file
		gelfEntry["_line"] = line
	}
	for key, value := range entry.Data {
		if !protectedFields[key] {
			key = "_" + key
		}
		gelfEntry[key] = value
	}
	message, err := json.Marshal(gelfEntry)
	return append(message, '\n'), err
}

func toTimestamp(t time.Time) float64 {
	nanosecond := float64(t.Nanosecond()) / 1e9
	seconds := float64(t.Unix())
	return seconds + nanosecond
}

func toSyslogLevel(level logrus.Level) syslog.Priority {
	syslog, ok := levelMap[level]
	if ok {
		return syslog
	}
	return DefaultLevel
}
