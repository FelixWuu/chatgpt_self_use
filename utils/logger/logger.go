package logger

import (
	"log"
	"os"
	"sync"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		Logger = log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	})
}

// Info 详情
func Info(args ...interface{}) {
	Logger.SetPrefix("[INFO]")
	Logger.Println(args...)
}

// Errorf 错误
func Errorf(format string, args ...interface{}) {
	Logger.SetPrefix("[ERROR]")
	Logger.Fatalf(format, args...)
}

// Warning 警告
func Warning(args ...interface{}) {
	Logger.SetPrefix("[WARNING]")
	Logger.Println(args...)
}

// DeBug debug
func DeBug(args ...interface{}) {
	Logger.SetPrefix("[DeBug]")
	Logger.Println(args...)
}
