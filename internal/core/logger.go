package core

import (
	"io"

	"github.com/clevergo/log"
	"github.com/clevergo/log/zapadapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
}

func NewLogger(cfg LogConfig) (log.Logger, func(), error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stderr"}
	config.ErrorOutputPaths = []string{"stderr"}
	sugar, err := config.Build()
	if err != nil {
		return nil, nil, err
	}

	undo := zap.RedirectStdLog(sugar)
	sugar = sugar.WithOptions(zap.AddCallerSkip(1))

	return zapadapter.New(sugar.Sugar()), func() {
		if err := sugar.Sync(); err != nil {
		}

		undo()
	}, nil
}

type loggerWriter struct {
	logger log.Logger
}

func (w *loggerWriter) Write(p []byte) (int, error) {
	w.logger.Info(string(p))
	return len(p), nil
}

func LoggerWriter(logger log.Logger) io.Writer {
	return &loggerWriter{logger: logger}
}
