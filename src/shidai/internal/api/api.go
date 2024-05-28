package api

import (
	"io"
	"net/http"

	"github.com/kiracore/sekin/src/shidai/internal/commands"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"github.com/nxadm/tail"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Serve() {
	router := gin.Default()
	router.POST("/api/execute", commands.ExecuteCommandHandler)
	router.GET("/logs/shidai", streamLogs(types.ShidaiLogPath))
	router.GET("/logs/sekai", streamLogs(types.SekaiLogPath))
	router.GET("/logs/interx", streamLogs(types.InterxLogPath))
	if err := router.Run(":8282"); err != nil {
		zap.L().Error("failed to run router...")
	}
}

func streamLogs(logFilePath string) gin.HandlerFunc {
	zap.L().Debug("streamLogs was called...")
	return func(c *gin.Context) {
		t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
		if err != nil {
			zap.L().Error("Unable to tail log file", zap.String("path", logFilePath), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to tail log file"})
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Flush() // Ensure the headers are sent immediately

		c.Stream(func(w io.Writer) bool {
			for line := range t.Lines {
				_, err := w.Write([]byte(line.Text + "\n"))
				if err != nil {
					zap.L().Error("Error writing to stream", zap.Error(err))
					return false // stop streaming on error
				}
				c.Writer.Flush()
			}
			return true // continue streaming
		})
	}
}
