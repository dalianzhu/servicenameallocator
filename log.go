package servicenameallocator

import "log"

var logger LogInterface = new(DefaultLogger)

func SetLogger(l LogInterface) {
	logger = l
}

type LogInterface interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type DefaultLogger struct{}

func (d *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func (d *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
