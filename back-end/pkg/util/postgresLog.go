package util

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

func NewPostgresLogger() (logger.Interface, error) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/postgres.log",
		MaxSize:    1,    // Max size in megabytes before log is rotated
		MaxBackups: 3,    // Max number of old log files to retain
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress/Archive old log files
	}

	LoggerConfig := logger.Config{
		SlowThreshold: 0,
		LogLevel:      logger.Info,
		Colorful:      true,
	}

	return logger.New(log.New(lumberjackLogger, "", log.LstdFlags), LoggerConfig), nil
}