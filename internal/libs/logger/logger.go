package logger

import (
	"fmt"
	"homedy/config"
	"os"
	"strings"
	"time"
)

func New(cfg ...Config) *Logger {
	c := Config{
		Output:     os.Stdout,
		Level:      DEBUG,
		TimeFormat: config.LOGGER_TIME_FORMAT,
		NoColor:    false,
	}
	if len(cfg) > 0 {
		if cfg[0].Output != nil {
			c.Output = cfg[0].Output
		}
		c.Level = cfg[0].Level
		if cfg[0].TimeFormat != "" {
			c.TimeFormat = cfg[0].TimeFormat
		}
		c.NoColor = cfg[0].NoColor
		c.Prefix = cfg[0].Prefix
	}
	return &Logger{config: c}
}

func (l *Logger) colorize(color, text string) string {
	if l.config.NoColor {
		return text
	}
	return color + text + colorReset
}

func (l *Logger) format(level Level, msg string, fields map[string]interface{}) string {
	now := time.Now().Format(l.config.TimeFormat)

	// Time
	timeStr := l.colorize(colorDim, now)

	// Level badge
	badge := l.colorize(level.color()+colorBold, level.badge())

	// Prefix
	prefix := ""
	if l.config.Prefix != "" {
		prefix = l.colorize(colorMagenta+colorBold, "["+l.config.Prefix+"] ")
	}

	// Message
	msgStr := l.colorize(colorWhite, msg)

	// Fields
	fieldStr := ""
	if len(fields) > 0 {
		parts := []string{}
		for k, v := range fields {
			key := l.colorize(colorCyan, k)
			val := l.colorize(colorYellow, fmt.Sprintf("%v", v))
			parts = append(parts, fmt.Sprintf("%s=%s", key, val))
		}
		fieldStr = "  " + strings.Join(parts, " ")
	}

	return fmt.Sprintf("%s %s %s%s%s\n", timeStr, badge, prefix, msgStr, fieldStr)
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	if level < l.config.Level {
		return
	}
	fmt.Fprint(l.config.Output, l.format(level, msg, fields))
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields...)
	l.log(DEBUG, msg, f)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields...)
	l.log(INFO, msg, f)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields...)
	l.log(WARN, msg, f)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields...)
	l.log(ERROR, msg, f)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	f := mergeFields(fields...)
	l.log(FATAL, msg, f)
	os.Exit(1)
}

// With returns a new logger with a prefix
func (l *Logger) With(prefix string) *Logger {
	cfg := l.config
	cfg.Prefix = prefix
	return &Logger{config: cfg}
}

// SetLevel changes the log level
func (l *Logger) SetLevel(level Level) {
	l.config.Level = level
}

// Package-level shortcuts using Default logger
func Debug(msg string, fields ...map[string]interface{}) { Default.Debug(msg, fields...) }
func Info(msg string, fields ...map[string]interface{})  { Default.Info(msg, fields...) }
func Warn(msg string, fields ...map[string]interface{})  { Default.Warn(msg, fields...) }
func Error(msg string, fields ...map[string]interface{}) { Default.Error(msg, fields...) }
func Fatal(msg string, fields ...map[string]interface{}) { Default.Fatal(msg, fields...) }
func With(prefix string) *Logger                         { return Default.With(prefix) }

// Fields builds a multi-key field map
func Fields(kv ...interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	for i := 0; i+1 < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		m[key] = kv[i+1]
	}
	return m
}

func mergeFields(fields ...map[string]interface{}) map[string]interface{} {
	merged := map[string]interface{}{}
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}
