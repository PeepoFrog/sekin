package genesishandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpexecutor "github.com/kiracore/sekin/src/shidai/internal/http_executor"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"go.uber.org/zap"
)

var (
	log = logger.GetLogger()
)

// unwrapGenesis checks if the genesis data has a "genesis" wrapper and unwraps it if present
func unwrapGenesis(data []byte, log *zap.Logger) ([]byte, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis data: %w", err)
	}

	// Check if there's a "genesis" wrapper
	if genesisWrapper, exists := raw["genesis"]; exists {
		// If the genesis key exists and it's the only key, unwrap it
		if len(raw) == 1 {
			log.Debug("Found genesis wrapper, unwrapping")
			return genesisWrapper, nil
		}
	}

	// No wrapper found, return original data
	log.Debug("No genesis wrapper found, using original format")
	return data, nil
}

// normalizeGenesisFormat ensures genesis file has consistent structure
// by unwrapping any "genesis" wrapper if present
func normalizeGenesisFormat(data []byte, log *zap.Logger) ([]byte, error) {
	unwrapped, err := unwrapGenesis(data, log)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize genesis format: %w", err)
	}

	// Re-marshal to ensure consistent formatting
	var genesis map[string]interface{}
	if err := json.Unmarshal(unwrapped, &genesis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal unwrapped genesis: %w", err)
	}

	normalized, err := json.Marshal(genesis)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal normalized genesis: %w", err)
	}

	log.Debug("Genesis format normalized successfully")
	return normalized, nil
}

func GetVerifiedGenesisFile(ctx context.Context, ip, interxPort string) ([]byte, error) {
	log.Info("Starting to get the genesis file", zap.String("IP", ip), zap.String("interxPort", interxPort))

	// Get genesis file from Interx daemon
	genesisInterx, err := GetInterxGenesis(ctx, ip, interxPort)
	if err != nil {
		log.Error("Failed to get genesis file from interx", zap.String("IP", ip), zap.String("Port", interxPort), zap.Error(err))
		return nil, fmt.Errorf("failed to get interx genesis: %w", err)
	}
	log.Debug("Retrieved genesis file from interx")

	// Normalize genesis file (unwrap if needed)
	normalized, err := normalizeGenesisFormat(genesisInterx, log)
	if err != nil {
		log.Error("Failed to normalize genesis", zap.Error(err))
		return nil, fmt.Errorf("failed to normalize genesis: %w", err)
	}

	log.Info("Genesis file retrieved and normalized successfully")
	return normalized, nil
}


func GetInterxGenesis(ctx context.Context, ipAddress, interxPort string) ([]byte, error) {
	log.Info("Starting to get the Interx genesis", zap.String("IP", ipAddress), zap.String("port", interxPort))
	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	// Construct the URL for fetching the genesis data
	url := fmt.Sprintf("http://%s:%s/api/tendermint/genesis", ipAddress, interxPort)
	log.Debug("Constructed URL for fetching Interx genesis", zap.String("url", url))

	// Create an HTTP client and perform the request
	client := &http.Client{}
	body, err := httpexecutor.DoHttpQuery(ctx, client, url, "GET")
	if err != nil {
		log.Error("Failed to get Interx genesis", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch Interx genesis from %s:%s: %w", ipAddress, interxPort, err)
	}

	log.Info("Interx genesis data retrieved successfully")
	return body, nil
}

