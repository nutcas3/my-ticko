package log

import "errors"

type Fields map[string]interface{}

const (
	Debug = "debug"
	Info = "info"
	Warn = "warn"
	Error = "error"
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

//Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
	Color             bool
}

//NewLogger returns an instance of logger
func NewLogger(config *Configuration, loggerInstance int) (Logger, error) {
	if config == nil {
		config = &Configuration{
			EnableConsole:     true,
			ConsoleLevel:      "debug",
			ConsoleJSONFormat: false,
			EnableFile:        false,
			Color: true,
		}
	}

	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(*config)
		if err != nil {
			return nil, err
		}
		return logger, nil
	case InstanceLogrusLogger:
		logger, err := newLogrusLogger(*config)
		if err != nil {
			return nil, err
		}
		return logger, nil

	default:
		return nil, errInvalidLoggerInstance
	}
}

func NormalizeLogLevel(logLevel string) string {
	var nomalizedLogLevel string
	switch logLevel {
	case "info":
		nomalizedLogLevel = Info
	case "debug":
		nomalizedLogLevel = Debug
	case "warn":
		nomalizedLogLevel = Warn
	case "error":
		nomalizedLogLevel = Error
	case "fatal":
		nomalizedLogLevel = Fatal
	default:
		nomalizedLogLevel = Debug
	}
	return nomalizedLogLevel
}
