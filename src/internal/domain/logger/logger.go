package logger

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Err(err error, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}
