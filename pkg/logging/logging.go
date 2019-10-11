package logging

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var baseLogger *logrus.Logger
var fieldLogger *logrus.Entry

var initOnce sync.Once

// Init sets the package level logger().
func Init(fields logrus.Fields) {
	baseLogger = logrus.New()

	baseLogger.SetFormatter(&logrus.TextFormatter{})

	logLevel := os.Getenv("PG_LOG_LEVEL")
	if logLevel != "" {
		parsedLevel, err := logrus.ParseLevel(logLevel)
		if err == nil {
			baseLogger.SetLevel(parsedLevel)
		} else {
			baseLogger.Errorf("Error while setting log level: %v", err)
		}
	} else {
		baseLogger.SetLevel(logrus.InfoLevel)
	}

	fieldLogger = baseLogger.WithFields(fields)
}

// Handle cases where the package level logger was never initialized.
func logger() *logrus.Entry {
	initOnce.Do(func() {
		if fieldLogger == nil {
			Init(logrus.Fields{})
		}
	})

	return fieldLogger
}

/*
LogJSONWithName logs json with caption.
*/
func LogJSONWithName(name string, ptr interface{}) {
	content, _ := json.MarshalIndent(ptr, "", " ")

	Devf("%s: %s", name, content)
}

/*
LogJSON logs json to the console.
*/
func LogJSON(ptr interface{}) {
	content, _ := json.MarshalIndent(ptr, "", " ")

	Dev(string(content))
}

// LogXMLWithName logs XML with a caption.
func LogXMLWithName(name string, ptr interface{}) {
	content, _ := xml.MarshalIndent(ptr, "", " ")

	Devf("%s: %s", name, content)
}

/*
LogRequest logs a raw http request.
*/
func LogRequest(r *http.Request) {

	content, _ := httputil.DumpRequest(r, true)
	Dev(string(content))

}

/*
LogResponse logs a raw http request.
*/
func LogResponse(r *http.Response) {

	content, _ := httputil.DumpResponse(r, true)
	Dev(string(content))

}

// Everything after this is just boilerplate to invoke logrus

// SetFormatter wraps the logrus SetFormatter function
func SetFormatter(formatter logrus.Formatter) {
	baseLogger.SetFormatter(formatter)
}

// AddHook wraps the logrus AddHook function
func AddHook(hook logrus.Hook) {
	baseLogger.AddHook(hook)
}

// WithFields wraps the logrus WithFields function.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger().WithFields(fields)
}

// Debugf wraps the logrus Debugf function
func Debugf(format string, args ...interface{}) {
	logger().Debugf(format, args...)
}

// Infof wraps the logrus Infof function
func Infof(format string, args ...interface{}) {
	logger().Infof(format, args...)
}

// Printf wraps the logrus Printf function
func Printf(format string, args ...interface{}) {
	logger().Printf(format, args...)
}

// Warnf wraps the logrus Warnf function
func Warnf(format string, args ...interface{}) {
	logger().Warnf(format, args...)
}

// Warningf wraps the logrus Warningf function
func Warningf(format string, args ...interface{}) {
	logger().Warnf(format, args...)
}

// Errorf wraps the logrus Errorf function
func Errorf(format string, args ...interface{}) {
	logger().Errorf(format, args...)
}

// Fatalf wraps the logrus Fatalf function
func Fatalf(format string, args ...interface{}) {
	logger().Fatalf(format, args...)
}

// Panicf wraps the logrus Panicf function
func Panicf(format string, args ...interface{}) {
	logger().Panicf(format, args...)
}

// Tracef wraps the logrus Tracef function
func Tracef(format string, args ...interface{}) {
	logger().Tracef(format, args...)
}

// Debug wraps the logrus Debug function
func Debug(args ...interface{}) {
	logger().Debug(args...)
}

// Info wraps the logrus Info function
func Info(args ...interface{}) {
	logger().Info(args...)
}

// Print wraps the logrus Print function
func Print(args ...interface{}) {
	logger().Info(args...)
}

// Warn wraps the logrus Warn function
func Warn(args ...interface{}) {
	logger().Warn(args...)
}

// Warning wraps the logrus Warning function
func Warning(args ...interface{}) {
	logger().Warn(args...)
}

// Error wraps the logrus Error function
func Error(args ...interface{}) {
	logger().Error(args...)
}

// Fatal wraps the logrus Fatal function
func Fatal(args ...interface{}) {
	logger().Fatal(args...)
}

// Panic wraps the logrus Panic function
func Panic(args ...interface{}) {
	logger().Panic(args...)
}

// Trace wraps the logrus Trace function
func Trace(args ...interface{}) {
	logger().Trace(args...)
}

// Debugln wraps the logrus Debugln function
func Debugln(args ...interface{}) {
	logger().Debugln(args...)
}

// Infoln wraps the logrus Infoln function
func Infoln(args ...interface{}) {
	logger().Infoln(args...)
}

// Println wraps the logrus Println function
func Println(args ...interface{}) {
	logger().Println(args...)
}

// Warnln wraps the logrus Warnln function
func Warnln(args ...interface{}) {
	logger().Warnln(args...)
}

// Warningln wraps the logrus Warningln function
func Warningln(args ...interface{}) {
	logger().Warnln(args...)
}

// Errorln wraps the logrus Errorln function
func Errorln(args ...interface{}) {
	logger().Errorln(args...)
}

// Fatalln wraps the logrus Fatalln function
func Fatalln(args ...interface{}) {
	logger().Fatalln(args...)
}

// Panicln wraps the logrus Panicln function
func Panicln(args ...interface{}) {
	logger().Panicln(args...)
}

// Traceln wraps the logrus Panicln function
func Traceln(args ...interface{}) {
	logger().Traceln(args...)
}
