# logx

包装zap日志包, 以及封装sugar, 使日志包更易使用. 并引入日志切割.

## 快速开始
### 代码
```go
import (
    "github.com/yann1989/logx"
)
    logger := logx.NewLogger(
    logx.WithRecordFile(logx.NewRecordFileWriter("./test.log", 10, 1, false)),
    logx.WithStdout(),
    logx.WithLevel(zapcore.DebugLevel))
    defer logger.Sync()
    logger.Info("测试1")
    logger.Infof("测试: %d", 1)
    return
```
### 输出
```text
2023-01-13 10:44:19	INFO	test/log_test.go:19	测试1
2023-01-13 10:44:19	INFO	test/log_test.go:20	测试: 1
```