package main

import (
	"time"

	"./logger"

	"go.uber.org/zap"
)

func main() {
	logger.Info("log 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "http://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
}
