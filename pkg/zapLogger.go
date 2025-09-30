package logger

import (
	"go.uber.org/zap"
	"os"
)

var log *zap.Logger

func InitLogger() {
	var err error
	switch os.Getenv("ENV") {
	case "dev", "development":
		log, err = zap.NewDevelopment()
	default:
		log, err = zap.NewProduction()
	}
	if err != nil {
		panic("cannot init zap logger: " + err.Error())
	}
}

func Sync() {
	_ = log.Sync()
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}
