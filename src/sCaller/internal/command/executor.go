package command

import (
	"fmt"
	"os/exec"
	"strings"
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
	if cmdArgs.Trace != "" {
		cmd.Args = append(cmd.Args, "--trace")
	}
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func SekaiAddGenesisAccCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(SekaiAddGenesisAcc)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'add-genesis-account'")
	}
	cmd := exec.Command(ExecPath, "add-genesis-account", cmdArgs.Address, strings.Join(cmdArgs.Coins, ","), "--home", cmdArgs.Home, "--keyring-backend", cmdArgs.Keyring, "--log_format", cmdArgs.LogFmt, "--log-level", cmdArgs.LogLvl)
	if cmdArgs.Trace != "" {
		cmd.Args = append(cmd.Args, "--trace")
	}
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func SekaiGentxClaimCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(SekaiGentxClaim)
	if !ok {

		return "", fmt.Errorf("invalid arguments for 'gentx-claim'")
	}
	cmd := exec.Command(ExecPath, "gentx-claim", cmdArgs.Address, "--keyring-backend", cmdArgs.Keyring, "--moniker", cmdArgs.Moniker, "--pubkey", cmdArgs.PubKey, "--home", cmdArgs.Home, "--log_format", cmdArgs.LogFmt, "--log-level", cmdArgs.LogLvl)
	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}
	output, err := cmd.CombinedOutput()

	return string(output), err
}
