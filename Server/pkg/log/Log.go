package log

type Information map[string]interface{}

type Log interface {
	Info(msg string, information Information)
	Warn(msg string, information Information)
	Error(msg string, err error, information Information)
}

var Logger Log = CreateLogger()
