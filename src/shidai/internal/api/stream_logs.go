package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nxadm/tail"
)

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
