// +build !prod

package logging

import (
	"encoding/json"
	"strings"
)

const devLogFormat = "[DEV LOG] "

// Dev is a log level that will only show up when running code compiled for
// development. When code is compiled for a release, dev will be a dummy
// function which does nothing.
func Dev(values ...interface{}) {
	arr := make([]interface{}, 0)
	arr = append(arr, devLogFormat)
	arr = append(arr, values...)
	logger().Warn(arr...)
}

// Devf is similar to dev, except that it accepts a format string.
func Devf(format string, values ...interface{}) {
	var b strings.Builder
	b.WriteString(devLogFormat)
	b.WriteString(format)
	logger().Warnf(b.String(), values...)
}

// DevJSON logs a struct as JSON. Similar to Dev, it will be compiled out in
// releases.
func DevJSON(value interface{}) {
	arr := make([]interface{}, 0)
	arr = append(arr, devLogFormat)
	content, _ := json.Marshal(value)
	arr = append(arr, string(content))
	logger().Warn(arr...)
}

// TempLog is kept for backwards compatibility. It passes through to Dev.
func TempLog(values ...interface{}) {
	Dev(values...)
}

// TempLogJSON is kept for backwards compatibility. It passes through to
// DevJSON.
func TempLogJSON(value interface{}) {
	DevJSON(value)
}
