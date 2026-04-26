package initialize

import (

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	logger, err := zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic("zap 初始化失败: " + err.Error())
	}

	zap.ReplaceGlobals(logger)
 
}
