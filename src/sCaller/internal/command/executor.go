package command

import (
	"fmt"
	"os/exec"
)

func SekaiInitCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(SekaiInit)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'init'")
	}
	cmd := exec.Command(ExecPath, "init", "--chain-id", cmdArgs.ChainID,
		fmt.Sprintf("--overwrite=%v", cmdArgs.Overwrite),
		"--log_format", cmdArgs.LogFmt,
		"--log-level", cmdArgs.LogLvl,
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func SekaiVersionCmd(args interface{}) (string, error) {
	cmd := exec.Command(ExecPath, "version")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func SekaidKeysAddCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(SekaidKeysAdd)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'init'")
	}

	cmd := exec.Command(ExecPath, "keys", "add", cmdArgs.KeyName,
		"--keyring-backend", cmdArgs.KeyringBackend,
		fmt.Sprintf("--recover=%v", cmdArgs.Recover),
		"--home", cmdArgs.Home,
		"--log_format", cmdArgs.LogFmt,
		"--log-level", cmdArgs.LogLvl,
	)

	output, err := cmd.CombinedOutput()
	return string(output), err
}
