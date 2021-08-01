package logger

import (
	"log"
	"os"
)

func NewWarnLogger() *log.Logger {
	return log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
}
