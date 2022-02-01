package utils

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger = nil
var atom = zap.NewAtomicLevel()

func SetLogLevel(level string) error {
	lvl := zapcore.InfoLevel // zero value
	err := lvl.Set(level)
	if err != nil {
		return err
	}
	atom.SetLevel(lvl)

	return nil
}

func Logger() *zap.Logger {
	if log == nil {
		InitLogger("console").Panic("Logger not yet initialized")
	}
	return log
}

func InitLogger(logEncoding string) *zap.Logger {
	cfg := zap.Config{
		Encoding:         logEncoding,
		Level:            atom,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			NameKey:      "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("could not create logger (%s)", err))
	}

	log = logger.Named("gowaku")
	return log
}
