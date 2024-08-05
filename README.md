# zdpgo_zap
基于zap二次开发，主要用于学习和研究

## 用法
### 基本使用
```go
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

```

## 版本历史
### v0.1.0
- 基础代码