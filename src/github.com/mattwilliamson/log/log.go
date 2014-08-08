package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

var G *Log

// Log allows logging with the date, file and line number to stdout
// Level will filter only equally or more important severities
type Log struct {
	*os.File
	sync.Mutex
	Level
}

// Add date, line no and file, if possible
func Format(level, format string, v ...interface{}) string {
	_, file, line, ok := runtime.Caller(3)
	file = path.Base(file)
	date := time.Now()
	datestr := date.Format(time.RFC3339)
	if ok {
		f := fmt.Sprintf("%v:%d", file, line)
		format = fmt.Sprintf("%v %-8v %-16v %v", datestr, level, f, format)
	} else {
		format = fmt.Sprintf("%v %-8v %v", datestr, level, format)
	}

	return fmt.Sprintf(format, v...)
}

// Log logs a message to stdout if the Level is less than the filter level
func (w *Log) Log(l Level, format string, v ...interface{}) (err error) {
	// Lock mutex
	w.Lock()
	defer w.Unlock()

	lname := LevelNames[l]

	// Add date and such to log line
	mf := Format(lname, format, v...)

	// Log to stdout
	_, err = w.File.WriteString(mf + "\n")

	if err != nil {
		return err
	}

	return err
}

// Emerg logs a message with severity LOG_EMERG, ignoring the severity
// passed to New.
func (w *Log) Emergency(m string, v ...interface{}) (err error) {
	return w.Log(LOG_EMERG, m, v...)
}

// Alert logs a message with severity LOG_ALERT, ignoring the severity
// passed to New.
func (w *Log) Alert(m string, v ...interface{}) (err error) {
	return w.Log(LOG_ALERT, m, v...)
}

// Crit logs a message with severity LOG_CRIT, ignoring the severity
// passed to New.
func (w *Log) Critical(m string, v ...interface{}) (err error) {
	return w.Log(LOG_CRIT, m, v...)
}

// Err logs a message with severity LOG_ERR, ignoring the severity
// passed to New.
func (w *Log) Error(m string, v ...interface{}) (err error) {
	return w.Log(LOG_ERR, m, v...)
}

// Wanring logs a message with severity LOG_WARNING, ignoring the
// severity passed to New.
func (w *Log) Warning(m string, v ...interface{}) (err error) {
	return w.Log(LOG_WARNING, m, v...)
}

// Notice logs a message with severity LOG_NOTICE, ignoring the
// severity passed to New.
func (w *Log) Notice(m string, v ...interface{}) (err error) {
	return w.Log(LOG_NOTICE, m, v...)
}

// Info logs a message with severity LOG_INFO, ignoring the severity
// passed to New.
func (w *Log) Info(m string, v ...interface{}) (err error) {
	return w.Log(LOG_INFO, m, v...)
}

// Debug logs a message with severity LOG_DEBUG, ignoring the severity
// passed to New.
func (w *Log) Debug(m string, v ...interface{}) (err error) {
	return w.Log(LOG_DEBUG, m, v...)
}

// Makes a new Log for logging to stdout
func New() (w *Log, err error) {
	l := &Log{
		File:   os.Stdout,
		Level:  LOG_INFO,
	}
	return l, err
}

// Set global default logger
func init() {
	G, _ = New()
}