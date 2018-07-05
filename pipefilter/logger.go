package pipefilter

// Logger for pipeline
type Logger interface {
	Info(v ...interface{})
}
