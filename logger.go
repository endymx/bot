package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger interface {
	Info(args ...any)
	Infof(template string, args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Debug(args ...any)
	Debugf(template string, args ...any)
	Panic(args ...any)
	Fatal(args ...any)
}

func SimpleLogger() *zap.SugaredLogger {
	syncConsole := zapcore.AddSync(os.Stderr)
	writeSyncer := zap.CombineWriteSyncers(syncConsole)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	defer func(SugarLogger *zap.SugaredLogger) {
		_ = SugarLogger.Sync()
	}(sugarLogger)

	return sugarLogger
}
