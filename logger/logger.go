package logger

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	gokitLogger "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/kolbis/corego/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents contract of logger
type Logger interface {
	Log(keyvals ...interface{}) error
}

// logger represents logger main struct , based on:
// kitLogger - gokit logger
// env - env. name where application runs
// loggerName - logger name(FileLogger , StdOutLogger etc.)
// minLevel - the minimal allowed level of logs
// dateUpdated - date of the latest Logger update , required for managing log files
type logger struct {
	kitLogger   gokitLogger.Logger
	env         string
	loggerName  string
	minLevel    AtomicLevelName
	dateUpdated time.Time
}

// Log implements Logger by calling logger.kitLogger.Log
// If logger.IsExpired return true -> logger will reload himselfe
func (log logger) Log(keyvals ...interface{}) error {
	if !log.IsLoggerDateValid() {
		log.Reload()
	}
	return log.kitLogger.Log(keyvals...)
}

// IsLoggerDateValid  return true if day of logger.dateCreated and time.Now are not equal
func (log logger) IsLoggerDateValid() bool {
	return log.dateUpdated.Day() == time.Now().UTC().Day()
}

// Reload function create new log.Logger object
func (log *logger) Reload() {
	dt := utils.DateTime{}
	dtNow := dt.Now()
	newLogger, err := buildLogger(dtNow)
	if err == nil {
		newLogger = gokitLogger.With(newLogger,
			"timestamp", gokitLogger.DefaultTimestampUTC,
			"caller", gokitLogger.Caller(8),
			"process", utils.ProcessName(),
			"loggerName", log.loggerName,
			"env", log.env,
		)
		newLogger = level.NewFilter(newLogger, toLevelOption(log.minLevel))
		log.kitLogger = newLogger
		log.dateUpdated = dtNow
	}
}

func buildLogger(dateNow time.Time) (log.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.LevelKey = ""
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.CallerKey = ""
	config.OutputPaths = append(config.OutputPaths, getOrCreatelogFilePath(dateNow))
	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}
	kitLogger := kitzap.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	return kitLogger, nil
}

// NewLogger create new logger which represents go-kit Logger
//logger object also contanins env , loggerName , minLevel and dateUpdated for a possibility to reload logger when it is expired
func NewLogger(env string, loggerName string, minLevel AtomicLevelName) (Logger, error) {
	dt := utils.DateTime{}
	dtNow := dt.Now()

	kitLogger, err := buildLogger(dtNow)

	if err != nil {
		return nil, err
	}

	kitLogger = gokitLogger.With(kitLogger,
		"timestamp", gokitLogger.DefaultTimestampUTC,
		"caller", gokitLogger.Caller(8),
		"process", utils.ProcessName(),
		"loggerName", loggerName,
		"env", env,
	)
	kitLogger = level.NewFilter(kitLogger, toLevelOption(minLevel))
	return &logger{
		kitLogger:   kitLogger,
		env:         env,
		loggerName:  loggerName,
		minLevel:    minLevel,
		dateUpdated: dtNow,
	}, nil
}

func getOrCreatelogFilePath(date time.Time) string {
	os.Mkdir("logs", os.ModePerm)
	fileName := "2006-01-02"
	filePath := "logs/" + date.Format(fileName)
	return filePath
}

func toLevelOption(l AtomicLevelName) level.Option {
	switch l {
	case DebugLogLevel:
		return level.AllowDebug()
	case InfoLogLevel:
		return level.AllowInfo()
	case WarnLogLevel:
		return level.AllowWarn()
	case ErrorLogLevel:
		return level.AllowError()
	case PanicLogLevel:
		return level.AllowError()
	default:
		return level.AllowAll()
	}
}

// AtomicLevelName represent name of specific log level
type AtomicLevelName string

const (
	// DebugLogLevel contains name of debug level
	DebugLogLevel AtomicLevelName = "DEBUG"
	// InfoLogLevel contains name of info level
	InfoLogLevel AtomicLevelName = "INFO"
	// WarnLogLevel contains name of warn level
	WarnLogLevel AtomicLevelName = "WARN"
	// ErrorLogLevel contains name of error level
	ErrorLogLevel AtomicLevelName = "ERROR"
	// PanicLogLevel contains name of panic level
	PanicLogLevel AtomicLevelName = "PANIC"
)
