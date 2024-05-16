package util

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

var APILog = &lumberjack.Logger{
	Filename:   "./logs/API.log",
	MaxSize:    1,    // Max size in megabytes before log is rotated
	MaxBackups: 3,    // Max number of old log files to retain
	MaxAge:     28,   // Max number of days to retain old log files
	Compress:   true, // Compress/Archive old log files
}