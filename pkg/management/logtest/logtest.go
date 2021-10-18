/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2021 EnterpriseDB Corporation.
*/

// Package logtest contains the testing utils for the logging subsystem of PGK
package logtest

import (
	"github.com/go-logr/logr"

	"github.com/EnterpriseDB/cloud-native-postgresql/pkg/management/log"
)

// LogLevel is the type representing a set of log levels
type LogLevel string

const (
	// LogLevelError is the error log level
	LogLevelError = LogLevel("ERROR")

	// LogLevelDebug is the debug log level
	LogLevelDebug = LogLevel("DEBUG")

	// LogLevelTrace is the error log level
	LogLevelTrace = LogLevel("TRACE")

	// LogLevelInfo is the error log level
	LogLevelInfo = LogLevel("INFO")
)

// LogRecord represents a log message
type LogRecord struct {
	LoggerName string
	Level      LogLevel
	Message    string
	Error      error
	Attributes map[string]interface{}
}

// NewRecord create a new log record
func NewRecord(name string, level LogLevel, msg string, err error, keysAndValues ...interface{}) *LogRecord {
	result := &LogRecord{
		LoggerName: name,
		Level:      level,
		Message:    msg,
		Error:      err,
		Attributes: make(map[string]interface{}),
	}
	result.WithValues(keysAndValues...)
	return result
}

// WithValues reads a set of keys and values, using them as attributes
// of the log record
func (record *LogRecord) WithValues(keysAndValues ...interface{}) {
	if len(keysAndValues)%2 != 0 {
		panic("key and values set is not even")
	}

	for idx := 0; idx < len(keysAndValues); idx += 2 {
		record.Attributes[keysAndValues[idx].(string)] = keysAndValues[idx+1]
	}
}

// SpyLogger is an implementation of the Logger interface that keeps track
// of the passed log entries
type SpyLogger struct {
	// The following attributes are referred to the current context

	Name       string
	Attributes map[string]interface{}

	// The following attributes represent the event sink

	Records   []LogRecord
	EventSink *SpyLogger
}

// NewSpy creates a new logger interface which will collect every log message sent
func NewSpy() *SpyLogger {
	result := &SpyLogger{Name: ""}
	result.EventSink = result
	return result
}

// AddRecord adds a log record inside the spy
func (s *SpyLogger) AddRecord(record *LogRecord) {
	s.EventSink.Records = append(s.EventSink.Records, *record)
}

// GetLogger implements the log.Logger interface
func (s SpyLogger) GetLogger() logr.Logger {
	return nil
}

// Enabled implements the log.Logger interface
func (s SpyLogger) Enabled() bool {
	return true
}

// Error implements the log.Logger interface
func (s *SpyLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	s.AddRecord(NewRecord(s.Name, LogLevelError, msg, err, keysAndValues...))
}

// Info implements the log.Logger interface
func (s SpyLogger) Info(msg string, keysAndValues ...interface{}) {
	s.AddRecord(NewRecord(s.Name, LogLevelInfo, msg, nil, keysAndValues...))
}

// Debug implements the log.Logger interface
func (s SpyLogger) Debug(msg string, keysAndValues ...interface{}) {
	s.AddRecord(NewRecord(s.Name, LogLevelDebug, msg, nil, keysAndValues...))
}

// Trace implements the log.Logger interface
func (s SpyLogger) Trace(msg string, keysAndValues ...interface{}) {
	s.AddRecord(NewRecord(s.Name, LogLevelTrace, msg, nil, keysAndValues...))
}

// WithValues implements the log.Logger interface
func (s SpyLogger) WithValues(keysAndValues ...interface{}) log.Logger {
	result := &SpyLogger{
		Name:      s.Name,
		EventSink: &s,
	}

	result.Attributes = make(map[string]interface{})
	for key, value := range s.Attributes {
		result.Attributes[key] = value
	}
	for idx := 0; idx < len(keysAndValues); idx += 2 {
		result.Attributes[keysAndValues[idx].(string)] = keysAndValues[idx+1]
	}

	return result
}

// WithName implements the log.Logger interface
func (s SpyLogger) WithName(name string) log.Logger {
	return &SpyLogger{
		Name:      name,
		EventSink: &s,
	}
}

// WithCaller implements the log.Logger interface
func (s SpyLogger) WithCaller() log.Logger {
	return &s
}