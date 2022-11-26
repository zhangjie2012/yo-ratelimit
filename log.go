package ratelimit

import "fmt"

// Logger basic logger
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// XLog a "Logger" basic implement by "fmt"
type XLog struct{}

func NewXLog() Logger {
	return &XLog{}
}

func (XLog) Debugf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (XLog) Infof(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (XLog) Errorf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
