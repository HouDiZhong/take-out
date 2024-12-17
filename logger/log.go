package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type ILog interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
}

type LLogger struct {
	logger *logrus.Logger
}

type LogEmailHook struct {
}

// Levels 需要监控的日志等级，只有命中列表中的日志等级才会触发Hook
func (l *LogEmailHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

// Fire 触发钩子函数，本实例为触发后发送邮件报警。
func (l *LogEmailHook) Fire(entry *logrus.Entry) error {
	// 触发loggerHook函数调用
	fmt.Println("触发loggerHook函数调用")
	return nil
}

func NewLogger(level string, filePath string) ILog {
	parseLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err.Error())
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile " + filePath)
		panic(err)
	}
	log := &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),         // 文件 + 控制台输出
		Level: parseLevel,                           // Debug日志等级
		Hooks: make(map[logrus.Level][]logrus.Hook), // 初始化Hook Map,否则导致Hook添加过程中的空指针引用。
		Formatter: &logrus.TextFormatter{ // 文本格式输出
			FullTimestamp:   true,                  // 展示日期
			TimestampFormat: "2006-01-02 15:04:05", //日期格式
			ForceColors:     false,                 // 颜色日志
		},
	}
	log.AddHook(&LogEmailHook{})
	log.Infof("日志开启成功")
	return &LLogger{logger: log}
}
func (l *LLogger) Debug(args ...any) {
	l.logger.Debug(args...)
}

func (l *LLogger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *LLogger) Warn(args ...any) {
	l.logger.Warn(args...)
}

func (l *LLogger) Error(args ...any) {
	l.logger.Error(args...)
}

func (l *LLogger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}
