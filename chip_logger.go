package chip

import (
	"log"
	"os"
)

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Warn(args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Panicf(template string, args ...any)
	Fatal(args ...any)
}

var logger Logger = &DefaultLogger{}

func setLogger(l Logger) {
	logger = l
}

type DefaultLogger struct{}

func (l *DefaultLogger) Debug(args ...any) {
	log.Println(args...)
}

func (l *DefaultLogger) Info(args ...any) {
	log.Println(args...)
}

func (l *DefaultLogger) Infof(template string, args ...any) {
	log.Printf(template, args...)
}

func (l *DefaultLogger) Warn(args ...any) {
	log.Println(args...)
}

func (l *DefaultLogger) Error(args ...any) {
	log.Println(args...)
}

func (l *DefaultLogger) Errorf(template string, args ...any) {
	log.Printf(template, args...)
}

func (l *DefaultLogger) Panicf(template string, args ...any) {
	log.Panicf(template, args...)
}

func (l *DefaultLogger) Fatal(args ...any) {
	log.Println(args...)
	os.Exit(1)
}
