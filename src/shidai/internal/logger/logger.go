package logger

import (
	"fmt"
	"log/syslog"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {

	// File path for backup logging
	logFilePath := "/syslog-data/syslog-ng/logs/shidai_backup.log"

	logFileFolder := filepath.Dir(logFilePath)

	if _, err := os.Stat(logFileFolder); os.IsNotExist(err) {
		err = os.MkdirAll(logFileFolder, 0644)
		if err != nil {
			panic(fmt.Sprintf("Unable to create logs folder %s. %s", logFileFolder, err.Error()))
		}
	}

	// Create the log file if it does not exist, or open it in append mode if it does
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Unable to create/open the log file: " + err.Error())
	}

	// Setting up network syslog writer
	syslogServer := "10.1.0.2:514" // Adjust as necessary
	syslogWriter, err := syslog.Dial("udp", syslogServer, syslog.LOG_LOCAL0, "shidai")
	if err != nil {
		panic("Failed to dial syslog: " + err.Error())
	}

	// Setting up encoders for each output
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	plaintextEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create cores for each output
	fileCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(logFile), zap.NewAtomicLevelAt(zapcore.InfoLevel))
	syslogCore := zapcore.NewCore(plaintextEncoder, zapcore.AddSync(syslogWriter), zap.NewAtomicLevelAt(zapcore.DebugLevel))

	// Combine cores
	combinedCore := zapcore.NewTee(fileCore, syslogCore)

	// Create the logger with the combined cores
	log = zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func GetLogger() *zap.Logger { return log }
