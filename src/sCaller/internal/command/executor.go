package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

func SekaiInitCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiInit)
	re := regexp.MustCompile(`\s+`)

	if !ok {
		return "", fmt.Errorf("invalid arguments for 'init'")
	}

	cmd := exec.Command(ExecPath, "init",
		"--home", cmdArgs.Home,
		"--chain-id", cmdArgs.ChainID,
		fmt.Sprintf("%v", re.ReplaceAllString(cmdArgs.Moniker, "_")),
		"--log_level", cmdArgs.LogLvl,
		"--log_format", cmdArgs.LogFmt,
	)

	if cmdArgs.Overwrite {
		cmd.Args = append(cmd.Args, "--overwrite")
	}

	log.Printf("DEBUG: SekaiInitCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiVersionCmd(interface{}) (string, error) {
	cmd := exec.Command(ExecPath, "version")
	log.Printf("DEBUG: SekaiVersionCmd: cmd: %v", cmd)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func SekaidKeysAddCmd(args interface{}) (string, error) {
	log.Printf("DEBUG: SekaidKeysAddCmd: in args: %v", args)
	cmdArgs, ok := args.(*SekaidKeysAdd)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'keys-add'")
	}

	cmd := exec.Command(ExecPath, "keys", "add", cmdArgs.Address,
		"--keyring-backend", cmdArgs.Keyring,
		"--home", cmdArgs.Home,
		"--log_format", cmdArgs.LogFmt,
		"--log_level", cmdArgs.LogLvl,
	)

	if cmdArgs.Output != "" {
		cmd.Args = append(cmd.Args, "--output", cmdArgs.Output)
	}
	if cmdArgs.Recover {
		cmd.Args = append(cmd.Args, "--recover")
	}
	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}

	log.Printf("DEBUG: SekaidKeysAddCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiAddGenesisAccCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiAddGenesisAcc)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'add-genesis-account'")
	}

	cmd := exec.Command(ExecPath, "add-genesis-account", cmdArgs.Address, strings.Join(cmdArgs.Coins, ","), "--home", cmdArgs.Home, "--keyring-backend", cmdArgs.Keyring, "--log_format", cmdArgs.LogFmt, "--log_level", cmdArgs.LogLvl)
	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}

	log.Printf("DEBUG: SekaiAddGenesisAccCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiGentxClaimCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiGentxClaim)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'gentx-claim'")
	}
	cmd := exec.Command(
		ExecPath, "gentx-claim", cmdArgs.Address,
		"--keyring-backend", cmdArgs.Keyring,
		"--moniker", fmt.Sprintf("%q", cmdArgs.Moniker),
		"--pubkey", cmdArgs.PubKey,
		"--home", cmdArgs.Home,
		"--log_format", cmdArgs.LogFmt,
		"--log_level", cmdArgs.LogLvl)

	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}
	log.Printf("DEBUG: SekaiGentxClaimCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))

	return string(output), err
}

func SekaidStartCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaidStart)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'start'")
	}

	argv := []string{"sekaid", "start", "--home", cmdArgs.Home}
	env := os.Environ()
	log.Printf("DEBUG: SekaidStartCmd: cmd args: %v", fmt.Sprintln(ExecPath, argv, env))
	err := syscall.Exec(ExecPath, argv, env)

	return "", err
}
