package main

import (
	zap "github.com/zhangdapeng520/zdpgo_zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	url := "http://github.com/zhangdapeng520"
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
