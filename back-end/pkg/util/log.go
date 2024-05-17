package util

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

// API log
var APILog = &lumberjack.Logger{
	Filename:   "./logs/API.log",
	MaxSize:    1,    // Max size in megabytes before log is rotated
	MaxBackups: 3,    // Max number of old log files to retain
	MaxAge:     28,   // Max number of days to retain old log files
	Compress:   true, // Compress/Archive old log files
}

// Crrate postgres sql logger
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
		Colorful:      false,
	}

	return logger.New(log.New(lumberjackLogger, "", log.LstdFlags), LoggerConfig), nil
}

// Create grpc logger
func GRPCLog() *log.Logger {
    lumberjackLogger := &lumberjack.Logger{
        Filename:   "./logs/grpc.log",
        MaxSize:    1, // megabytes
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true, // disabled by default
    }
    
    return log.New(lumberjackLogger, "", log.LstdFlags)
}

func KafkaLogger() *log.Logger {
    lumberjackLogger := &lumberjack.Logger{
        Filename:   "./logs/kafka.log",
        MaxSize:    10, // megabytes
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true, // disabled by default
    }

    return log.New(lumberjackLogger, "", log.LstdFlags)
}