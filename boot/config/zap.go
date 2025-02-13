package config

import (
	"go-tasks/boot/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var Logger *zap.Logger
var logLevel string = "info"

type ZapLogger struct {
	level     zapcore.Level
	file      *os.File
	mutex     sync.Mutex
	encoder   zapcore.Encoder
	directory string
	filename  string
}

func createLogDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create log directory:", err)
		}
	}
}

func InitZapLogger() {
	createLogDir(global.LogDir)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller", // 包含行号
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		//EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	lv, _ := zapcore.ParseLevel(logLevel)

	logger := &ZapLogger{
		level:     lv,
		directory: global.LogDir,
		encoder:   zapcore.NewConsoleEncoder(encoderConfig),
		//encoder:  zapcore.NewJSONEncoder(encoderConfig),
		filename: time.Now().Format("2006-01-02"),
	}

	core := zapcore.NewCore(
		logger.encoder,
		logger,
		lv,
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(Logger)
}

func (z *ZapLogger) Sync() error {
	if z.file != nil {
		return z.file.Sync()
	}
	return nil
}

func (z *ZapLogger) Write(p []byte) (n int, err error) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	// 检查文件是否存在或是否需要轮转
	date := time.Now().Format("2006-01-02")
	if z.file == nil || date != z.filename {
		if z.file != nil {
			err := z.file.Close()
			if err != nil {
				return 0, err
			}
		}
		z.filename = date
		filePath := filepath.Join(z.directory, z.filename+".log")
		z.file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
	}

	p = append(p, '\n')

	return z.file.Write(p)
}
