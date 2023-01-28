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
	sep     string
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
	if logger.writers == nil || len(logger.writers) == 0 {
		logger.writers = []io.Writer{os.Stdout}
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
	encoderConfig.ConsoleSeparator = l.sep
	return l.encoder(encoderConfig)
}

func (l *Logger) getSyncer() zapcore.WriteSyncer {
	var syncers []zapcore.WriteSyncer
	for i := 0; i < len(l.writers); i++ {
		syncers = append(syncers, zapcore.AddSync(l.writers[i]))
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
		sugar:  l.sugar.Named(name),
	}
}
