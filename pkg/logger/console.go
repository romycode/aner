package logger

import "log"

type StdoutLogger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func NewStdoutLogger(il *log.Logger, wl *log.Logger, el *log.Logger) *StdoutLogger {
	return &StdoutLogger{
		info:  il,
		warn:  wl,
		error: el,
	}
}

func (s StdoutLogger) Info(m interface{}) {
	s.info.Println(m)
}

func (s StdoutLogger) Warn(m interface{}) {
	s.warn.Println(m)
}

func (s StdoutLogger) Error(m error) {
	s.error.Fatalln(m)
}
