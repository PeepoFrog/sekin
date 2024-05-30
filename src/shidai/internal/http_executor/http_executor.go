package httpexecutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kiracore/sekin/src/shidai/internal/utils"
	"go.uber.org/zap"
)

type CommandRequest struct {
	Command string      `json:"command"`
	Args    interface{} `json:"args"`
}

// Executes command for iCaller and sCaller
func ExecuteCallerCommand(address, port, method string, commandRequest CommandRequest) ([]byte, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("<%v> port is not valid", port)
	}
	check := utils.ValidatePort(p)
	if !check {
		return nil, fmt.Errorf("<%v> port is not valid", port)
	}
	// Convert your struct to JSON
	jsonData, err := json.Marshal(commandRequest)
	if err != nil {
		zap.L().Debug("Error marshaling JSON:", zap.Error(err))
		return nil, err
	}

	// Create a new request
	req, err := http.NewRequest(method, fmt.Sprintf("http://%v:%v/api/execute", address, port), bytes.NewBuffer(jsonData))
	if err != nil {
		zap.L().Debug("Error creating request:", zap.Error(err))
		return nil, err
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Debug("Error making request:", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Debug("Error reading response body:", zap.Error(err))
		return nil, err
	}

	// fmt.Println("Response:", string(body))
	return body, nil
}

func DoHttpQuery(ctx context.Context, client *http.Client, url, method string) ([]byte, error) {
	const timeoutQuery = time.Second * 10
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		zap.L().Debug("ERROR: Failed to create request:", zap.Error(err))
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.L().Debug("ERROR: Failed to send request: %s", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Debug("ERROR: Failed to read response body: %s", zap.Error(err))
		return nil, err
	}

	return body, nil
}
