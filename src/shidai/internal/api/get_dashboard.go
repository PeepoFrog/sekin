package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sm "github.com/kiracore/sekin/src/shidai/internal/subscriptionmanager"
	"go.uber.org/zap"
)

func backgroundUpdate(manager *sm.SubscriptionManager) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Assuming updateDashboardCache modifies the dashboardCache with new data
		updateDashboardCache(someNewData, manager)

		// Now persist the updated data
		lock.Lock()
		bytes, err := json.Marshal(dashboardCache)
		if err != nil {
			log.Error("Failed to marshal dashboard cache", zap.Error(err))
		} else if err := os.WriteFile(dataFile, bytes, 0640); err != nil {
			log.Error("Failed to write dashboard cache to file", zap.String("file", dataFile), zap.Error(err))
		} else {
			log.Info("Dashboard cache updated and saved to file", zap.String("file", dataFile))
		}
		lock.Unlock()
	}
}

func getDashboard(c *gin.Context) {
	lock.RLock()
	defer lock.RUnlock()

	bytes, err := os.ReadFile(dataFile)
	if err != nil {
		log.Error("Failed to read dashboard cache", zap.String("file", dataFile), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch dashboard"})
		return
	}
	c.Data(http.StatusOK, "application/json", bytes)
}
