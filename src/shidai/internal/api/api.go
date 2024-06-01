package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiracore/sekin/src/shidai/internal/commands"
	"github.com/kiracore/sekin/src/shidai/internal/docker"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"github.com/kiracore/sekin/src/shidai/internal/utils"
	"github.com/nxadm/tail"
	"go.uber.org/zap"
)

func Serve() {
	log := logger.GetLogger()

	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/api/execute", commands.ExecuteCommandHandler)
	router.GET("/logs/shidai", streamLogs(types.ShidaiLogPath))
	router.GET("/logs/sekai", streamLogs(types.SekaiLogPath))
	router.GET("/logs/interx", streamLogs(types.InterxLogPath))
	router.GET("/status", infraStatus())

	if err := router.Run(":8282"); err != nil {
		log.Error("Failed to start the server", zap.Error(err))
	}
}

func streamLogs(logFilePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to tail log file"})
			return
		}

		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Flush()

		c.Stream(func(w io.Writer) bool {
			for line := range t.Lines {
				_, err := w.Write([]byte(line.Text + "\n"))
				if err != nil {
					return false
				}
				c.Writer.Flush()
			}
			return true
		})
	}
}

func infraStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			SekaiVersion, InterxVersion, SyslogVersion, ShidaiVersion []byte
			err                                                       error
		)
		cm, err := docker.NewContainerManager()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		ctx := c.Request.Context()

		if SekaiVersion, err = cm.ExecInContainer(ctx, "sekin-sekai-1", []string{"/sekaid", "version"}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Sekai version"})
			return
		}

		if InterxVersion, err = cm.ExecInContainer(ctx, "sekin-interx-1", []string{"/interxd", "version"}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Interx version"})
			return
		}
		if SyslogVersion, err = cm.ExecInContainer(ctx, "sekin-syslog-ng-1", []string{"syslog-ng", "--version"}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Syslog version"})
			return
		}

		if ShidaiVersion, err = cm.ExecInContainer(ctx, "sekin-shidai-1", []string{"/shidai", "version"}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Shidai version"})
			return
		}

		response := types.StatusResponse{
			Sekai:  types.AppInfo{Version: string(SekaiVersion), Infra: utils.CheckInfra(types.SekaiFiles)},
			Interx: types.AppInfo{Version: string(InterxVersion), Infra: utils.CheckInfra(types.InterxFiles)},
			Syslog: types.AppInfo{Version: string(SyslogVersion), Infra: utils.CheckInfra(types.SyslogFiles)},
			Shidai: types.AppInfo{Version: string(ShidaiVersion), Infra: true},
		}
		c.JSON(http.StatusOK, response)
	}
}
