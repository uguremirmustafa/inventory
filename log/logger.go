package logging

import (
	"fmt"
	"os"
	"time"
)

// Logger interface defines the methods for logging.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// FileLogger implements the Logger interface and logs messages to a file.
type FileLogger struct {
	filename string
}

// NewFileLogger creates a new FileLogger instance.
func NewFileLogger(filename string) *FileLogger {
	return &FileLogger{filename: filename}
}

// log writes a log message to the file.
func (l *FileLogger) log(level, format string, args ...interface{}) {
	f, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
		return
	}
	defer f.Close()

	message := fmt.Sprintf("[%s] [%s] %s\n", level, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	if _, err := f.WriteString(message); err != nil {
		fmt.Fprintf(os.Stderr, "error writing to log file: %v\n", err)
	}
}

// Debugf logs a debug message.
func (l *FileLogger) Debugf(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

// Infof logs an info message.
func (l *FileLogger) Infof(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Warnf logs a warning message.
func (l *FileLogger) Warnf(format string, args ...interface{}) {
	l.log("WARNING", format, args...)
}

// Errorf logs an error message.
func (l *FileLogger) Errorf(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}
