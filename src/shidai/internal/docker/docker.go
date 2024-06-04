package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"go.uber.org/zap"
)

// ContainerManager manages Docker containers.
type ContainerManager struct {
	Cli *client.Client
}

var log = logger.GetLogger() // Initialize the logger instance at the package level

// NewContainerManager creates a new instance of ContainerManager.
func NewContainerManager() (*ContainerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error("Failed to create Docker client", zap.Error(err))
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	log.Info("Docker client initialized successfully")
	return &ContainerManager{Cli: cli}, nil
}

func GetAccAddress(ctx context.Context, cm *ContainerManager, containerID string) (string, error) {
	command := []string{"/sekaid", "keys", "show", "validator", "--home", "/sekai", "--keyring-backend", "test", "--output", "json"}
	output, err := cm.ExecInContainer(ctx, containerID, command)
	if err != nil {
		log.Error("Failed to execute command in container", zap.String("containerID", containerID), zap.Error(err))
		return "", fmt.Errorf("failed to execute command in container %s: %w", containerID, err)
	}

	var result struct {
		Address string `json:"address"`
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		log.Error("Failed to parse JSON output", zap.Error(err))
		return "", fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return result.Address, nil
}

func GetRoleID(ctx context.Context, cm *ContainerManager, containerID, address string) ([]string, error) {
	command := []string{
		"/sekaid", "q", "customgov", "roles", address,
		"--node", "tcp://sekai.local:26657",
		"--output", "json",
	}
	output, err := cm.ExecInContainer(ctx, containerID, command)
	if err != nil {
		log.Error("Failed to execute command in container", zap.String("containerID", containerID), zap.Error(err))
		return nil, fmt.Errorf("failed to execute command in container %s: %w", containerID, err)
	}

	var result struct {
		RoleIDs []string `json:"roleIds"`
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		log.Error("Failed to parse JSON output", zap.Error(err))
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	return result.RoleIDs, nil
}

// ExecInContainer executes a command inside a specified container and returns the output.
func (cm *ContainerManager) ExecInContainer(ctx context.Context, containerID string, command []string) ([]byte, error) {
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
	}
	execCreateResponse, err := cm.Cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		log.Error("Failed to create container exec instance", zap.String("containerID", containerID), zap.Error(err))
		return nil, fmt.Errorf("failed to create container exec instance for container %s: %w", containerID, err)
	}

	execAttachConfig := types.ExecStartCheck{}
	resp, err := cm.Cli.ContainerExecAttach(ctx, execCreateResponse.ID, execAttachConfig)
	if err != nil {
		log.Error("Failed to attach to container exec instance", zap.String("execID", execCreateResponse.ID), zap.Error(err))
		return nil, fmt.Errorf("failed to attach to container exec instance %s: %w", execCreateResponse.ID, err)
	}
	defer resp.Close()

	var outBuf, errBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
	if err != nil {
		log.Error("Failed to copy output from container exec", zap.Error(err))
		return nil, fmt.Errorf("failed to copy output from container exec: %w", err)
	}

	if len(errBuf.Bytes()) > 0 {
		log.Warn("Standard error output from container exec", zap.ByteString("stderr", errBuf.Bytes()))
	}

	output := outBuf.Bytes()
	log.Info("Command executed successfully", zap.String("output", string(output)))
	return output, nil
}
