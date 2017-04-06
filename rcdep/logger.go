package rcdep

type Logger interface {
	Printf(format string, v ...interface{})
}
