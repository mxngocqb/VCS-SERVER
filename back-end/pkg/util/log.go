package util

import "gopkg.in/natefinch/lumberjack.v2"

var LogConfig = &lumberjack.Logger{
	Filename:   "./logs/server.log",
	MaxSize:    5,    // megabytes
	MaxBackups: 10,   // max number of backup files
	MaxAge:     30,   // days
	Compress:   true, // compress rolled files
}
