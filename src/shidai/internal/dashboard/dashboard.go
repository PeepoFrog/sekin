package dashboard

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"go.uber.org/zap"
)

var (
	log *zap.Logger = logger.GetLogger()
)

func fetchAndUpdateDashboard() {
	ticker := time.NewTicker(10 * time.Minute) // Update every 10 minutes
	defer ticker.Stop()

	for range ticker.C {
		newData, err := fetchDataFromEndpoint("https://api.yourdomain.com/dashboard")
		if err != nil {
			log.Error("Failed to fetch data", zap.Error(err))
			continue
		}

		if err := UpdateDashboardData(*newData); err != nil {
			log.Error("Failed to update dashboard data", zap.Error(err))
		}
	}
}

func fetchDataFromEndpoint(url string) (*types.Dashboard, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newData types.Dashboard
	if err := json.Unmarshal(body, &newData); err != nil {
		return nil, err
	}

	return &newData, nil
}

func UpdateDashboardData(newData types.Dashboard) error {
	lock.Lock()
	defer lock.Unlock()

	// Update the cache; this part depends on how newData should be merged with existing data.
	dashboardCache = newData

	// Marshal the updated dashboard data into JSON
	bytes, err := json.Marshal(dashboardCache)
	if err != nil {
		log.Error("Failed to marshal dashboard cache", zap.Error(err))
		return err
	}

	// Attempt to write the JSON data to the file, with retries
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if err := os.WriteFile(dataFile, bytes, 0644); err != nil {
			log.Error("Failed to write dashboard cache to file", zap.String("file", dataFile), zap.Error(err))
			if i < maxRetries-1 {
				log.Info("Retrying file write", zap.Int("attempt", i+2))
				time.Sleep(time.Second) // Simple backoff strategy
				continue
			}
			return err
		}
		break
	}

	log.Info("Dashboard cache updated and saved successfully", zap.String("file", dataFile))
	return nil
}
