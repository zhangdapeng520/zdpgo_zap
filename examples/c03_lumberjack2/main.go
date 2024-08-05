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
