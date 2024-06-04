package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiracore/sekin/src/shidai/internal/docker"
	"github.com/kiracore/sekin/src/shidai/internal/types"
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
		ChainID             string   `json:"chain_id"`
		NodeID              string   `json:"node_id"`
		RoleIDs             []string `json:"role_ids"`
		SeatClaimAvailable  bool     `json:"seat_claim_available"`
		GenesisChecksum     string   `json:"gen_sha256"`
	}

	DashboardPointer struct {
		Data *Dashboard
		mu   sync.RWMutex
	}
)

var dashboardPointer = NewDashboardPointer()

func NewDashboardPointer() *DashboardPointer {
	return &DashboardPointer{
		Data: &Dashboard{
			Date:                "",
			ValidatorStatus:     "",
			Blocks:              "",
			Top:                 "",
			Streak:              "",
			Mischance:           "",
			MischanceConfidence: "",
			StartHeight:         "",
			LastProducedBlock:   "",
			ProducedBlocks:      "",
			Moniker:             "",
			ValidatorAddress:    "",
			ChainID:             "",
			NodeID:              "",
			RoleIDs:             []string{""},
			SeatClaimAvailable:  false,
			GenesisChecksum:     "",
		},
	}
}

func getDashboardHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		dashboardPointer.mu.RLock()
		clone := dashboardPointer.Data.DeepCopy() // Use the read lock to allow concurrent reads
		defer dashboardPointer.mu.RUnlock()

		// Instead of reading from file, directly marshal the in-memory dashboard
		dashboardJSON, err := json.Marshal(clone)
		if err != nil {
			log.Error("Failed to marshal dashboard data", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal dashboard data"})
			return
		}

		c.Data(http.StatusOK, "application/json", dashboardJSON)
	}
}

func backgroundUpdate() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Debug("Starting background update")

		if err := updateDashboard(); err != nil {
			log.Error("Failed to update dashboard", zap.Error(err))
			continue
		}

		if err := saveDashboardToFile(); err != nil {
			log.Error("Failed to save dashboard to file", zap.String("file", types.DashboardPath), zap.Error(err))
		}
	}
}

func saveDashboardToFile() error {
	dashboardPointer.mu.Lock()
	defer dashboardPointer.mu.Unlock()

	writer, err := os.OpenFile(types.DashboardPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, types.FilePermRW)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer writer.Close()

	bufWriter := bufio.NewWriter(writer)
	if err := json.NewEncoder(bufWriter).Encode(dashboardPointer.Data); err != nil {
		return fmt.Errorf("failed to marshal and write dashboard: %w", err)
	}
	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush write buffer: %w", err)
	}

	log.Info("Dashboard pointer updated and saved to file", zap.String("file", types.DashboardPath))
	return nil
}

func updateDashboard() error {
	ctx := context.Background()
	cm, err := docker.NewContainerManager()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	dashboardData := dashboardPointer.GetDashboardData()
	dashboardUpdates := make(chan *Dashboard, 10) // Buffer size based on expected concurrency
	done := make(chan error, 10)

	wg.Add(3) //Increase with qty of fetches

	go func() {
		defer wg.Done()
		fetchDateNow(dashboardUpdates, done)
	}()

	go func() {
		defer wg.Done()
		fetchAccAddressFromSekaidBin(ctx, cm, types.SEKAI_CONTAINER_ID, dashboardUpdates, done)
	}()

	go func() {
		defer wg.Done()
		fetchRoleIDsFromSekaidBin(ctx, cm, types.SEKAI_CONTAINER_ID, dashboardData.ValidatorAddress, dashboardUpdates, done)
	}()

	var errors []error
	for i := 0; i < 3; i++ { // Adjust this based on the number of goroutines
		if err := <-done; err != nil {
			log.Warn("Failed to fetch some data", zap.Error(err))
			errors = append(errors, err)
		}
	}

	go func() {
		wg.Wait()
		close(dashboardUpdates)
		close(done)
	}()
	if len(errors) > 0 {
		log.Warn("Fetching ecnounted errors", zap.String("times", strconv.Itoa(len(errors))))
		for _, err := range errors {
			log.Error("Detail[FETCH]: ", zap.Error(err))
		}
	}

	for update := range dashboardUpdates {
		applyUpdate(dashboardPointer, update) // Apply updates to the central dashboardPointer
	}

	return nil
}

// applyUpdate safely applies updates to the shared dashboardPointer.
func applyUpdate(pointer *DashboardPointer, update *Dashboard) {
	// Lock the pointer fo r writing.
	pointer.mu.Lock()
	defer pointer.mu.Unlock()

	// Check and apply each field if the incoming update is not the default "".
	// This example assumes "" is the default state for string fields and checks accordingly.

	if update.Date != "" {
		pointer.Data.Date = update.Date
	}
	pointer.Data.ValidatorStatus = update.ValidatorStatus
	pointer.Data.Blocks = update.Blocks
	pointer.Data.Top = update.Top
	pointer.Data.Streak = update.Streak
	pointer.Data.Mischance = update.Mischance
	pointer.Data.MischanceConfidence = update.MischanceConfidence
	pointer.Data.StartHeight = update.StartHeight
	pointer.Data.LastProducedBlock = update.LastProducedBlock
	pointer.Data.ProducedBlocks = update.ProducedBlocks
	pointer.Data.Moniker = update.Moniker
	if update.ValidatorAddress != "" {
		pointer.Data.ValidatorAddress = update.ValidatorAddress
	}
	pointer.Data.ValidatorAddress = update.ValidatorAddress
	pointer.Data.ChainID = update.ChainID
	pointer.Data.NodeID = update.NodeID
	if len(update.RoleIDs) > 0 {
		pointer.Data.RoleIDs = update.RoleIDs
	}
	pointer.Data.SeatClaimAvailable = update.SeatClaimAvailable
	pointer.Data.GenesisChecksum = update.GenesisChecksum
}

func fetchRoleIDsFromSekaidBin(ctx context.Context, cm *docker.ContainerManager, containerID, adr string, updates chan<- *Dashboard, done chan<- error) {
	defer func() { done <- nil }()
	log.Debug("Fetching roleID from sekai container")
	if containerID == "" || adr == "" {
		done <- fmt.Errorf("containerID or address can't be empty")
		return
	}
	command := []string{"/sekaid", "q", "customgov", "roles", adr, "--node", types.SEKAI_RPC_LADDR, "--output", "json"}
	ouput, err := cm.ExecInContainer(ctx, containerID, command)
	if err != nil {
		done <- fmt.Errorf("failed to execute command in container %s: %w", containerID, err)
		return
	}
	var result struct {
		RoleIDs []string `json:"roleids"`
	}
	if err = json.Unmarshal(ouput, &result); err != nil {
		done <- fmt.Errorf("failed to parse JSON output: %w", err)
		return

	}
	log.Info("Sending update to dashboardUpdates")
	updates <- &Dashboard{RoleIDs: result.RoleIDs}
}

func fetchAccAddressFromSekaidBin(ctx context.Context, cm *docker.ContainerManager, containerID string, updates chan<- *Dashboard, done chan<- error) {
	defer func() { done <- nil }()
	log.Debug("Fetching address from sekai container")
	if containerID == "" {
		done <- fmt.Errorf("containerID can't be empty")
		return
	}
	command := []string{"/sekaid", "keys", "show", "validator", "--home", "/sekai", "--keyring-backend", "test", "--output", "json"}
	output, err := cm.ExecInContainer(ctx, containerID, command)
	if err != nil {
		done <- fmt.Errorf("failed to execute command in container %s: %w", containerID, err)
		return
	}

	var result struct {
		Address string `json:"address"`
	}
	if err = json.Unmarshal(output, &result); err != nil {
		done <- fmt.Errorf("faile d to parse JSON output: %w", err)
		return
	}
	log.Info("Sending update to dashboardUpdates")
	updates <- &Dashboard{ValidatorAddress: result.Address}
}

//	func fetchDummy(updates chan<- *Dashboard, done chan<- error) {
//		defer func() { done <- nil }()
//		done <- fmt.Errorf("error")
//		updates <- &Dashboard{Key: Value}
//	}

func fetchDateNow(updates chan<- *Dashboard, done chan<- error) {
	defer func() { done <- nil }()
	log.Debug("Fetching date ")
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05")
	log.Info("Sending update to dashboardUpdates")
	updates <- &Dashboard{Date: formattedDate}

}

// DeepCopy creates a deep copy of the Dashboard.
func (d *Dashboard) DeepCopy() *Dashboard {
	clone := *d // Copy all primitive fields

	// Manually copy slices to ensure they are independent
	if d.RoleIDs != nil {
		clone.RoleIDs = make([]string, len(d.RoleIDs))
		copy(clone.RoleIDs, d.RoleIDs)
	}

	return &clone
}

// GetDashboardData retrieves a deep copy of the Dashboard data safely.
func (dp *DashboardPointer) GetDashboardData() *Dashboard {
	dp.mu.RLock()         // Acquire a read lock
	defer dp.mu.RUnlock() // Ensure the lock is released when the function exits

	return dp.Data.DeepCopy()
}
