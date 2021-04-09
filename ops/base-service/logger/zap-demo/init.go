package main

import (
	"os"

	"./logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	fileSyncer    lumberjack.Logger
	encoderConfig zapcore.EncoderConfig
)

func init() {
	initFileSyncer()
	initenEoderConfig()

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                                 // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&fileSyncer)), // 打印到控制台和文件
		atomicLevel, // 日志级别配置
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "testservice"))

	logger.Init(zap.New(core, caller, development, filed))
}

func initFileSyncer() {
	fileSyncer = lumberjack.Logger{
		Filename:   "./logs/test-service.log", // 日志文件路径
		MaxSize:    1,                         // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                        // 日志文件最多保存多少个备份
		MaxAge:     7,                         // 文件最多保存多少天
		Compress:   true,                      // 是否压缩
	}
}

func initenEoderConfig() {
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
}
