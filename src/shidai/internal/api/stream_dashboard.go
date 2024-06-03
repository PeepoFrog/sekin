package api

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	sm "github.com/kiracore/sekin/src/shidai/internal/subscriptionmanager"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"go.uber.org/zap"
)

var (
	dataFile       = "/sekai/dashboard_data.json"
	lock           = &sync.RWMutex{}
	dashboardCache types.Dashboard
)

func streamDashboard(manager *sm.SubscriptionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Setting SSE stream headers")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		updateChan := make(chan bool)
		defer close(updateChan)

		manager.Subscribe(updateChan)
		defer manager.Unsubscribe(updateChan)

		// Send the current state immediately upon connection
		lock.RLock()
		c.SSEvent("message", dashboardCache)
		lock.RUnlock()
		c.Writer.Flush()
		log.Info("Sent initial SSEvent", zap.String("event", "update_dashboard"))
		log.Debug("Initial dashboard cache sent as SSEvent", zap.Any("cache", dashboardCache))

		for {
			select {
			case <-updateChan:
				lock.RLock()
				c.SSEvent("message", dashboardCache)
				lock.RUnlock()
				c.Writer.Flush()
				log.Info("Sent SSEvent", zap.String("event", "update_dashboard"))
				log.Debug("Dashboard cache sent as SSEvent", zap.Any("cache", dashboardCache))
			case <-c.Writer.CloseNotify():
				log.Debug("Client has disconnected")
				return
			}
		}
	}
}

func initCache() error {
	err := createDefaultCacheFile()
	if err != nil {
		log.Fatal("Failed to create cache file for Dashboard")
	}
	log.Debug("Initializing dashboard cache from file", zap.String("file", dataFile))
	bytes, err := os.ReadFile(dataFile)
	if err != nil {
		log.Error("Failed to read cache file", zap.String("file", dataFile), zap.Error(err))
		return err
	}
	if err := json.Unmarshal(bytes, &dashboardCache); err != nil {
		log.Error("Failed to unmarshal dashboard cache", zap.Error(err))
		return err
	}
	log.Info("Dashboard cache initialized successfully")
	return nil
}

func persistData() {
	log.Debug("Starting periodic data persistence")
	ticker := time.NewTicker(10 * time.Second) // Adjust interval based on expected update frequency
	defer ticker.Stop()

	for range ticker.C {
		lock.Lock()
		bytes, err := json.Marshal(dashboardCache)
		if err != nil {
			log.Error("Failed to marshal dashboard cache", zap.Error(err))
			lock.Unlock()
			continue
		}
		if err := os.WriteFile(dataFile, bytes, 0640); err != nil {
			log.Error("Failed to write dashboard cache to file", zap.String("file", dataFile), zap.Error(err))
		} else {
			log.Info("Dashboard cache saved to file", zap.String("file", dataFile))
		}
		lock.Unlock()
	}
}

func createDefaultCacheFile() error {
	log.Debug("Creating default cache file", zap.String("file", dataFile))

	// Create a default Dashboard object with the current date and time
	defaultDashboard := types.NewDashboard()
	bytes, err := json.MarshalIndent(defaultDashboard, "", "  ")
	if err != nil {
		log.Error("Failed to marshal default dashboard data", zap.Error(err))
		return err
	}

	// Write the default data to the file
	if err := os.WriteFile(dataFile, bytes, 0644); err != nil {
		log.Error("Failed to write default dashboard cache to file", zap.String("file", dataFile), zap.Error(err))
		return err
	}

	log.Info("Default dashboard cache file created successfully with current date and time", zap.String("file", dataFile))
	return nil
}

func updateDashboardCache(newData interface{}, manager *sm.SubscriptionManager) {
	lock.Lock()
	dashboardCache = types.NewDashboard()
	lock.Unlock()

	manager.NotifySubscribers()
}

// func updateCachePeriodically() {
// 	ticker := time.NewTicker(3 * time.Second) // Adjust based on how frequently you need updates
// 	defer ticker.Stop()
//
// 	for range ticker.C {
// 		newData, err := fetchDataFromSource()
// 		if err != nil {
// 			log.Error("Failed to fetch new data", zap.Error(err))
// 			continue
// 		}
//
// 		lock.Lock()
// 		updateDashboardCache(newData)
// 		lock.Unlock()
//
// 		log.Info("Dashboard cache updated with new data")
// 	}
// }
