package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger initializes the zap logger
func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Setting the time format

	// Create a file to write logs
	logFile, err := os.Create("shidai_logs.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file: %v", err))
	}

	// Write logs to the file in JSON format
	fileWriteSyncer := zapcore.AddSync(logFile)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		fileWriteSyncer,
		config.Level,
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(Log) // Replace the global logger, which can be accessed with zap.L()
}
