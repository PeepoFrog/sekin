package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiracore/sekin/src/shidai/internal/docker"
	"go.uber.org/zap"
)

type (
	Dashboard struct {
		Date                string   `json:"date"`
		ValidatorStatus     string   `json:"val_status"`
		Blocks              string   `json:"blocks"`
		Top                 string   `json:"top"`
		Streak              string   `json:"streak"`
		Mischance           string   `json:"mischance"`
		MischanceConfidence string   `json:"mischance_confidence"`
		StartHeight         string   `json:"start_height"`
		LastProducedBlock   string   `json:"last_present_block"`
		ProducedBlocks      string   `json:"produced_blocks_counter"`
		Moniker             string   `json:"moniker"`
		ValidatorAddress    string   `json:"address"`
		NodeID              string   `json:"node_id"`
		RoleID              []string `json:"role_ids"`
		GenesisChecksum     string   `json:"gen_sha256"`
		SeatClaimAvailable  bool     `json:"seat_claim_available"`
	}
)

var (
	dashboardCache *Dashboard
	lock           sync.RWMutex
	dataFile       = "/shidaid/dashboard_cache.json"
)

func NewDashboard() Dashboard {
	return Dashboard{
		Date:                getCurrentTimeString(), // Current date and time in YYYY-MM-DD HH:MM:SS format
		ValidatorStatus:     "Unknown",
		Blocks:              "Unknown",
		Top:                 "Unknown",
		Streak:              "Unknown",
		Mischance:           "Unknown",
		MischanceConfidence: "Unknown",
		StartHeight:         "Unknown",
		LastProducedBlock:   "Unknown",
		ProducedBlocks:      "Unknown",
		Moniker:             "Unknown",
		ValidatorAddress:    "Unknown",
		NodeID:              "Unknown",
		GenesisChecksum:     "Unknown",
		SeatClaimAvailable:  false,
	}
}

func backgroundUpdate() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Debug("Starting background update")

		// Update the dashboard data
		updateDashboard()

		// Persist the updated data to the file
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

func getDashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func updateDashboard() {
	cm, _ := docker.NewContainerManager()

	addr, err := docker.GetAccAddress(context.Background(), cm, "sekin-sekai-1")
	if err != nil {
		log.Error("Failed to get account address", zap.Error(err))
		addr = "Unknown"
	}

	id, err := docker.GetRoleID(context.Background(), cm, "sekin-sekai-1", addr)
	if err != nil {
		log.Error("Failed to get role ID", zap.Error(err))
		id = []string{"Unknown"}
	}

	resp, err := fetchValidatorData(addr)
	if err != nil {
		log.Error("Failed to fetch validator data", zap.Error(err))
		// Set relevant fields to "Unknown" if an error occurs
		dashboardCache = &Dashboard{
			Date:                getCurrentTimeString(),
			ValidatorAddress:    addr,
			NodeID:              "Unknown",
			RoleID:              id,
			ValidatorStatus:     "Unknown",
			Top:                 "Unknown",
			Streak:              "Unknown",
			Mischance:           "Unknown",
			MischanceConfidence: "Unknown",
			StartHeight:         "Unknown",
			LastProducedBlock:   "Unknown",
			ProducedBlocks:      "Unknown",
			Moniker:             "Unknown",
		}
	} else {
		dashboard, err := mapResponseToDashboard(resp)
		if err != nil {
			log.Error("Failed to map response to dashboard", zap.Error(err))
			// Set relevant fields to "Unknown" if an error occurs
			dashboardCache = &Dashboard{
				Date:                getCurrentTimeString(),
				ValidatorAddress:    addr,
				NodeID:              "Unknown",
				RoleID:              id,
				ValidatorStatus:     "Unknown",
				Top:                 "Unknown",
				Streak:              "Unknown",
				Mischance:           "Unknown",
				MischanceConfidence: "Unknown",
				StartHeight:         "Unknown",
				LastProducedBlock:   "Unknown",
				ProducedBlocks:      "Unknown",
				Moniker:             "Unknown",
			}
		} else {
			dashboardCache = dashboard
			dashboardCache.Date = getCurrentTimeString()
			dashboardCache.ValidatorAddress = addr
			dashboardCache.RoleID = id
			dashboardCache.GenesisChecksum = getFileSHA256(dataFile)
			dashboardCache.SeatClaimAvailable = containsTwo(dashboardCache.RoleID)

		}
	}
}

func fetchValidatorData(address string) ([]byte, error) {
	url := fmt.Sprintf("http://interx.local:11000/api/valopers?address=%s", address)
	log.Debug("Fetching validator data", zap.String("url", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func mapResponseToDashboard(jsonResponse []byte) (*Dashboard, error) {
	var response struct {
		Validators []struct {
			Top                   string `json:"top"`
			Address               string `json:"address"`
			Moniker               string `json:"moniker"`
			Status                string `json:"status"`
			Streak                string `json:"streak"`
			Mischance             string `json:"mischance"`
			MischanceConfidence   string `json:"mischance_confidence"`
			StartHeight           string `json:"start_height"`
			LastPresentBlock      string `json:"last_present_block"`
			ProducedBlocksCounter string `json:"produced_blocks_counter"`
		} `json:"validators"`
	}

	err := json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	if len(response.Validators) == 0 {
		return nil, fmt.Errorf("no validator found in the response")
	}

	validator := response.Validators[0]
	dashboard := &Dashboard{
		Date:                getCurrentTimeString(),
		ValidatorStatus:     validator.Status,
		Top:                 validator.Top,
		Streak:              validator.Streak,
		Mischance:           validator.Mischance,
		MischanceConfidence: validator.MischanceConfidence,
		StartHeight:         validator.StartHeight,
		LastProducedBlock:   validator.LastPresentBlock,
		ProducedBlocks:      validator.ProducedBlocksCounter,
		Moniker:             strings.Trim(validator.Moniker, "\""),
	}

	return dashboard, nil
}

func getFileSHA256(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("failed to open ", zap.String("file", filePath), zap.Error(err))
		return "Unknown"
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {

		log.Error("failed to get hash ", zap.String("file", filePath), zap.Error(err))
		return "Unknown"
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}
func containsTwo(arr interface{}) bool {
	// Check if the input is a slice of strings
	strArr, ok := arr.([]string)
	if !ok {
		return false
	}

	limit := len(strArr)
	if limit > 10 {
		limit = 10
	}

	for i := 0; i < limit; i++ {
		// Convert each element to a string and compare with "2"
		if fmt.Sprint(strArr[i]) == "2" {
			return true
		}
	}

	return false
}

func getCurrentTimeString() string {
	now := time.Now()
	formatted := now.Format("2006-01-02 15:04:05") // Format as "YYYY-MM-DD HH:mm:ss"
	return formatted
}
