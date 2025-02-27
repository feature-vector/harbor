package starter

import (
	"github.com/feature-vector/harbor/base/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZap() {
	if env.IsProduction() {
		infoConfig := zap.NewProductionEncoderConfig()
		infoConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		infoCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(infoConfig),
			zapcore.AddSync(&lumberjack.Logger{Filename: "logs/app.log"}),
			zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l == zapcore.WarnLevel || l == zapcore.InfoLevel
			}),
		)
		errorConfig := zap.NewProductionEncoderConfig()
		errorConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		errorConfig.EncodeCaller = zapcore.FullCallerEncoder
		errCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(errorConfig),
			zapcore.AddSync(&lumberjack.Logger{Filename: "logs/error.log"}),
			zap.ErrorLevel,
		)
		zap.ReplaceGlobals(zap.New(
			zapcore.NewTee(infoCore, errCore),
			zap.AddStacktrace(zap.ErrorLevel),
		))
	} else {
		devLogger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(devLogger)
	}
}
