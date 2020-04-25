package gologutil

import (
	"fmt"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	CONF_FILE       = "log.json"
	LOG_FILE        = "app.log"
	LOG_PORT        = ":8888"
	LOG_LEVEL       = "debug"
	LOG_MAX_SIZE    = 100 //100MB
	LOG_MAX_BACKUPS = 3
	LOG_MAX_AGE     = 30 //Days

)

var (
	logFile, logPort, logLevel  string
	logger                      *zap.Logger
	aLevel                      zap.AtomicLevel
	maxSize, maxBackups, maxAge int64
)

func init() {

	conf, confErr := ioutil.ReadFile(CONF_FILE)

	if confErr != nil || !gjson.ValidBytes(conf) {

		logFile = LOG_FILE
		logPort = LOG_PORT
		logLevel = LOG_LEVEL
		maxAge = LOG_MAX_AGE
		maxBackups = LOG_MAX_BACKUPS
		maxSize = LOG_MAX_SIZE

	} else {

		if gjson.GetBytes(conf, "logFile").String() != "" {
			logFile = gjson.GetBytes(conf, "logFile").String()
		} else {
			logFile = LOG_FILE
		}

		if gjson.GetBytes(conf, "logPort").String() != "" {
			logPort = ":" + gjson.GetBytes(conf, "logPort").String()
		} else {
			logPort = LOG_PORT
		}

		if gjson.GetBytes(conf, "logLevel").String() != "" {
			logLevel = strings.ToLower(gjson.GetBytes(conf, "logLevel").String())
		} else {
			logLevel = LOG_LEVEL
		}

		if gjson.GetBytes(conf, "maxSize").Uint() != 0 {
			maxSize = gjson.GetBytes(conf, "maxSize").Int()
		} else {
			maxSize = LOG_MAX_SIZE
		}

		if gjson.GetBytes(conf, "maxBackups").Uint() != 0 {
			maxBackups = gjson.GetBytes(conf, "maxBackups").Int()
		} else {
			maxBackups = LOG_MAX_BACKUPS
		}

		if gjson.GetBytes(conf, "maxAge").Uint() != 0 {
			maxAge = gjson.GetBytes(conf, "maxAge").Int()
		} else {
			maxAge = LOG_MAX_AGE
		}

	}

	switch logLevel {
	case "info":
		aLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		aLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "fatal":
		aLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		aLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	http.HandleFunc("/loglevel", aLevel.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(logPort, nil); err != nil {
			panic(err)
		}
	}()

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    int(maxSize), // megabytes
		MaxBackups: int(maxBackups),
		MaxAge:     int(maxAge), // days
	})

	encoder := zap.NewProductionEncoderConfig()
	//屏蔽调用位置
	encoder.CallerKey = ""
	//格式化时间显示方式
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//文件和控制台同时输出
	syncConfig := zapcore.NewMultiWriteSyncer(w, zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncConfig, aLevel)
	logger = zap.New(core)

}

func makeZapFiels(fields map[string]string) []zap.Field {
	var zapFileds []zap.Field
	for key, value := range fields {
		zapFileds = append(zapFileds, zap.String(key, value))
	}
	return zapFileds
}

func Debug(message string, key string, value string) {
	logger.Debug(message, zap.String(key, value))
}

func DebugWithFields(message string, fields map[string]string) {
	logger.Debug(message, makeZapFiels(fields)...)
}

func Info(message string, key string, value string) {
	logger.Info(message, zap.String(key, value))
}

func InfoWithFields(message string, fields map[string]string) {
	logger.Info(message, makeZapFiels(fields)...)
}

func Warn(message string, key string, value string) {
	logger.Warn(message, zap.String(key, value))
}

func WarnWithFields(message string, fields map[string]string) {
	logger.Warn(message, makeZapFiels(fields)...)
}

func Fatal(message string, key string, value string) {
	logger.Fatal(message, zap.String(key, value))
}

func FatalWithFields(message string, fields map[string]string) {
	logger.Fatal(message, makeZapFiels(fields)...)
}

func Sync() {
	logger.Sync()
}

func PrintConf() {
	fmt.Println("logFile:", logFile)
	fmt.Println("logLevel:", logLevel)
	fmt.Println("logPort:", logPort)
	fmt.Println("maxSize(MB):", maxSize)
	fmt.Println("maxBackups:", maxBackups)
	fmt.Println("maxAge(Days):", maxAge)
}
