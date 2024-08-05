package zdpgo_zap

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_zap/lumberjack"
	"github.com/zhangdapeng520/zdpgo_zap/zapcore"
	"os"
	"time"
)

type LumberjackLoggerConfig struct {
	Level      string `json:"level"`       // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string `json:"file_name"`   // FileName 日志文件位置
	MaxSize    int    `json:"max_size"`    // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int    `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int    `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
	IsConsole  bool   `json:"is_console"`  // 是否写入到控制台
}

func NewLumberjackLogger(config *LumberjackLoggerConfig) (logger *Logger, err error) {
	// 读取配置信息
	lc := LumberjackLoggerConfig{
		Level:      "debug",
		FileName:   fmt.Sprintf("./data/logs/%v.log", time.Now().Unix()),
		MaxSize:    99,
		MaxBackups: 33,
		MaxAge:     33,
		IsConsole:  true,
	}
	if config.Level != "" {
		lc.Level = config.Level
	}
	if config.FileName != "" {
		lc.FileName = config.FileName
	}
	if config.MaxSize != 0 {
		lc.MaxSize = config.MaxSize
	}
	if config.MaxBackups != 0 {
		lc.MaxBackups = config.MaxBackups
	}
	if config.MaxAge != 0 {
		lc.MaxAge = config.MaxAge
	}

	// 初始化日志对象
	return initLogger(lc)
}

// 负责设置 encoding 的日志格式
func getEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConfig := NewProductionEncoderConfig()

	// 设置每个日志条目使用的键。如果有任何键为空，则省略该条目的部分。
	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "time"
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 以 package/file:行 的格式 序列化调用程序，从完整路径中删除除最后一个目录外的所有目录。
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

// 负责日志写入的位置
func getLogWriter(filename string, maxsize, maxBackup, maxAge int, isConsole bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer。
	// 它试图变得智能：如果 io.Writer 的具体类型实现了 WriteSyncer，我们将使用现有的 Sync 方法。
	// 如果没有，我们将添加一个无操作同步。
	syncFile := zapcore.AddSync(lumberJackLogger) // 打印到文件
	if isConsole {
		syncConsole := zapcore.AddSync(os.Stderr) // 打印到控制台
		return zapcore.NewMultiWriteSyncer(syncFile, syncConsole)
	} else {
		return zapcore.AddSync(lumberJackLogger)
	}
}

// initLogger 初始化Logger
func initLogger(config LumberjackLoggerConfig) (logger *Logger, err error) {
	// 获取日志写入位置
	writeSyncer := getLogWriter(config.FileName, config.MaxSize, config.MaxBackups, config.MaxAge, config.IsConsole)
	// 获取日志编码格式
	encoder := getEncoder()

	// 获取日志最低等级，即>=该等级，才会被写入。
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(config.Level))
	if err != nil {
		return
	}

	// 创建一个将日志写入 WriteSyncer 的核心。
	core := zapcore.NewCore(encoder, writeSyncer, l)
	logger = New(core, AddCaller())

	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	ReplaceGlobals(logger)
	return
}
