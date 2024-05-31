package logger

import (
	"log/syslog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {

	// File path for backup logging
	logFilePath := "/syslog-data/syslog-ng/logs/backup.log"

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

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create two sync outputs for both writers
	combinedSync := zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile), zapcore.AddSync(syslogWriter))

	// Create a core that writes to both file and syslog
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		combinedSync,
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	)

	// Create the logger with the defined core
	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	//zap.ReplaceGlobals(Log)

}

func GetLogger() *zap.Logger { return log }

// InitLogger initializes the zap logger
// func InitLogger() {
// 	config := zap.NewProductionConfig()
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Setting the time format
//
// 	// Create a file to write logs
// 	logFile, err := os.Create("/syslog-data/syslog-ng/logs/shidai.log")
// 	if err != nil {
// 	}
//
// 	// Write logs to the file in JSON format
// 	fileWriteSyncer := zapcore.AddSync(logFile)
// 	core := zapcore.NewCore(
// 		zapcore.NewJSONEncoder(config.EncoderConfig),
// 		fileWriteSyncer,
// 		config.Level,
// 	)
//
// 	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
// }
