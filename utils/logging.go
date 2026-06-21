// Package utils provides utility functions for logging, directory creation,
// and other common operations used throughout the application.
package utils

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// LogIfError logs an error if the provided message is not nil.
func LogIfError(msg interface{}) {
	if msg == nil {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Errorf("%s: return not nil: %s", elements[len(elements)-1], msg)
}

// LogStart logs the start of a function
func LogStart() {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Debugf("%s: start", elements[len(elements)-1])

}

// LogEnd logs the end of a function
func LogEnd() {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Debugf("%s: end", elements[len(elements)-1])
}

// LogArgument logs a function argument with its name and value.
// Uses runtime introspection to automatically determine the calling function name.
func LogArgument(name, input interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf(" argument: %s(%T)\n", name, input))
	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf("    value: %#v\n", input))

}

// LogVariable logs a variable with its name and value.
// Uses runtime introspection to automatically determine the calling function name.
func LogVariable(name, input interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf(" variable: %s(%T)\n", name, input))
	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf("    value: %#v\n", input))

}

// Debugf logs a debug message with the name of the function that called it
func Debugf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Debugf("%s: %s", elements[len(elements)-1], msg)

}

// Infof logs an info message with the name of the function that called it
func Infof(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Infof("%s: %s", elements[len(elements)-1], msg)

}

// Errorf logs an error message with the name of the function that called it
func Errorf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Errorf("%s: %s", elements[len(elements)-1], msg)

}

// Warningf logs an warning message with the name of the function that called it
func Warningf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Warningf("%s: %s", elements[len(elements)-1], msg)

}

// Fatalf logs a fatal error message with the name of the function that called it
func Fatalf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("%s: %s", elements[len(elements)-1], msg)

}
