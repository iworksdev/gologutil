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
	ConfFile      = "log.json"
	LogFile       = "app.log"
	LogPort       = ":8888"
	LogLevel      = "debug"
	LogMaxSize    = 100 //100MB
	LogMaxBackups = 3
	LogMaxAge     = 30 //Days

)

var (
	logFile, logPort, logLevel  string
	logger                      *zap.Logger
	aLevel                      zap.AtomicLevel
	maxSize, maxBackups, maxAge int64
	consoleOutput               bool
)

func init() {

	conf, confErr := ioutil.ReadFile(ConfFile)

	if confErr != nil || !gjson.ValidBytes(conf) {

		logFile = LogFile
		logPort = LogPort
		logLevel = LogLevel
		maxAge = LogMaxAge
		maxBackups = LogMaxBackups
		maxSize = LogMaxSize
		consoleOutput = true

	} else {

		if gjson.GetBytes(conf, "logFile").Exists() {
			logFile = gjson.GetBytes(conf, "logFile").String()
		} else {
			logFile = LogFile
		}

		if gjson.GetBytes(conf, "logPort").Exists() {
			logPort = ":" + gjson.GetBytes(conf, "logPort").String()
		} else {
			logPort = LogPort
		}

		if gjson.GetBytes(conf, "logLevel").Exists() {
			logLevel = strings.ToLower(gjson.GetBytes(conf, "logLevel").String())
		} else {
			logLevel = LogLevel
		}

		if gjson.GetBytes(conf, "consoleOutput").Exists() {
			consoleOutput = gjson.GetBytes(conf, "consoleOutput").Bool()
		} else {
			consoleOutput = true
		}

		if gjson.GetBytes(conf, "maxSize").Exists() {
			maxSize = gjson.GetBytes(conf, "maxSize").Int()
		} else {
			maxSize = LogMaxSize
		}

		if gjson.GetBytes(conf, "maxBackups").Exists() {
			maxBackups = gjson.GetBytes(conf, "maxBackups").Int()
		} else {
			maxBackups = LogMaxBackups
		}

		if gjson.GetBytes(conf, "maxAge").Exists() {
			maxAge = gjson.GetBytes(conf, "maxAge").Int()
		} else {
			maxAge = LogMaxAge
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
	var syncConfig zapcore.WriteSyncer
	if consoleOutput {
		syncConfig = zapcore.NewMultiWriteSyncer(w, zapcore.AddSync(os.Stdout))
	} else {
		syncConfig = zapcore.WriteSyncer(w)
	}

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

func Error(message string, key string, value string) {
	logger.Error(message, zap.String(key, value))
}

func ErrorWithFields(message string, fields map[string]string) {
	logger.Error(message, makeZapFiels(fields)...)
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
	fmt.Println("consoleOutput:", consoleOutput)
}
