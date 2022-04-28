package internal

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
}

type NoOpLogger struct{}

func (NoOpLogger) Print(...any)          {}
func (NoOpLogger) Printf(string, ...any) {}
