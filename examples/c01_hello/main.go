package main

import (
	"github.com/zhangdapeng520/zdpgo_zap"
	"time"
)

func main() {
	logger, _ := zdpgo_zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	url := "http://github.com/zhangdapeng520"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
