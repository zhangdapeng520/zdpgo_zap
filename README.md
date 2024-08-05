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

### 实现日志切割功能
```go
package main

import (
	"fmt"
	zap "github.com/zhangdapeng520/zdpgo_zap"
	"net/http"
	"time"
)

func main() {
	logger, err := zap.NewLumberjackLogger(&zap.LumberjackLoggerConfig{
		Level:      "debug",
		FileName:   fmt.Sprintf("./log/%v.log", time.Now().Unix()),
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		IsConsole:  true,
	})
	if err != nil {
		panic(err)
	}
	// 调用内核的Sync方法，刷新所有缓冲的日志条目。
	// 应用程序应该注意在退出之前调用Sync。
	defer logger.Sync()
	simpleHttpGet(logger, "www.sogo.com")
	simpleHttpGet(logger, "http://www.sogo.com")
}

func simpleHttpGet(logger *zap.Logger, url string) {
	sugarLogger := logger.Sugar()
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

```

## 版本历史

### v0.1.0

- 基础代码

### v0.1.1

- 整合lumberjack