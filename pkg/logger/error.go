package logger

import (
	"log"
	"os"
)

func NewErrorLogger() *log.Logger {
	return log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
