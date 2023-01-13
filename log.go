package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Logger struct {
	*zap.Logger
	sugar   *zap.SugaredLogger
	l       zapcore.Level
	writers []io.Writer
	format  string
	encoder func(cfg zapcore.EncoderConfig) zapcore.Encoder
	hooks   []func(zapcore.Entry) error
}

func NewLogger(options ...Option) *Logger {
	logger := &Logger{
		l:       zapcore.DebugLevel,
		format:  "2006-01-02 15:04:05",
		encoder: zapcore.NewConsoleEncoder,
	}
	for _, option := range options {
		option(logger)
	}
	if len(logger.writers) == 0 {
		logger.writers = append(logger.writers, os.Stdout)
	}

	encoder := logger.getEncoder()
	syncer := logger.getSyncer()
	core := zapcore.NewCore(encoder, syncer, logger.l)
	logger.Logger = zap.New(core, zap.AddCaller())
	if len(logger.hooks) != 0 {
		logger.Logger = logger.Logger.WithOptions(zap.Hooks(logger.hooks...))
	}
	logger.sugar = logger.Sugar().WithOptions(zap.AddCallerSkip(1))
	return logger
}

func (l *Logger) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(l.format)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return l.encoder(encoderConfig)
}

func (l *Logger) getSyncer() zapcore.WriteSyncer {
	var syncers []zapcore.WriteSyncer
	for i := 0; i < len(l.writers); i++ {
		syncers = append(syncers, zapcore.AddSync(l.writers[i]))
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

func WithLevel(l zapcore.Level) Option {
	return func(logger *Logger) {
		logger.l = l
	}
}
