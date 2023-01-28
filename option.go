package logx

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Option func(logger *Logger)

//WithJson 是否输出json格式
func WithJson() Option {
	return func(logger *Logger) {
		logger.encoder = zapcore.NewJSONEncoder
	}
}

//WithTimeFormat 时间输出格式
func WithTimeFormat(format string) Option {
	return func(logger *Logger) {
		logger.format = format
	}
}

//WithStdout 是否标准输出
func WithStdout() Option {
	return func(logger *Logger) {
		logger.writers = append(logger.writers, os.Stdout)
	}
}

//WithRecordFile 是否输出文件
func WithRecordFile(writer io.Writer) Option {
	return func(logger *Logger) {
		logger.writers = append(logger.writers, writer)
	}
}

func WithHooks(hooks ...func(zapcore.Entry) error) Option {
	return func(logger *Logger) {
		logger.hooks = hooks
	}
}

func WithLevel(l zapcore.Level) Option {
	return func(logger *Logger) {
		logger.l = l
	}
}

func WithSeparator(sep string) Option {
	return func(logger *Logger) {
		logger.sep = sep
	}
}

//NewRecordFileWriter
//file 文件路径
//maxSize 单个文件大小 单位MB
//maxAge  文件存在时间(到期自动删除) 单位天
//compress 是否压缩
func NewRecordFileWriter(file string, maxSize, maxAge int, compress bool) io.Writer {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: 0,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
