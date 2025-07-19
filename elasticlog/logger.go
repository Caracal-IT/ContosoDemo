package elasticlog

import (
	"fmt"
	"os"
	"strings"
)

// LogLevel is the type for log levels.
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var logLevelStrings = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}

func (l LogLevel) String() string {
	if int(l) < 0 || int(l) >= len(logLevelStrings) {
		return "UNKNOWN"
	}
	return logLevelStrings[l]
}

// Logger logs to both console and Elasticsearch.
type Logger struct {
	Level LogLevel
	Index string
	User  string
	Pass  string
}

// NewLogger creates a new Logger.
func NewLogger(level LogLevel, index, user, pass string) *Logger {
	return &Logger{
		Level: level,
		Index: index,
		User:  user,
		Pass:  pass,
	}
}

// logInternal logs to console and Elasticsearch if level is enough.
func (l *Logger) logInternal(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.Level {
		return
	}
	// Console log
	consoleMsg := fmt.Sprintf("[%s] %s", level.String(), msg)
	if len(fields) > 0 {
		consoleMsg += " | " + fmt.Sprint(fields)
	}
	_, _ = fmt.Fprintln(os.Stdout, consoleMsg)

	// Elastic log
	LogToElastic(level.String(), msg, fields, l.Index, l.User, l.Pass)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.logInternal(DebugLevel, msg, fields)
}
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.logInternal(InfoLevel, msg, fields)
}
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.logInternal(WarnLevel, msg, fields)
}
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.logInternal(ErrorLevel, msg, fields)
}

// ParseLogLevel parses a string to LogLevel, defaults to InfoLevel.
func ParseLogLevel(s string) LogLevel {
	s = strings.ToUpper(s)
	for i, v := range logLevelStrings {
		if v == s {
			return LogLevel(i)
		}
	}
	return InfoLevel
}
