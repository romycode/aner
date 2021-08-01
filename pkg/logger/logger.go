package logger

type Logger interface {
	Info(m interface{})
	Warn(m interface{})
	Error(m error)
}
