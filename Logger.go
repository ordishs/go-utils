package utils

import "log"

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args)
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("WARN: "+format, args)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args)
}

func (l *defaultLogger) Fatalf(format string, args ...interface{}) {
	println("FATAL: "+format, args)
}
