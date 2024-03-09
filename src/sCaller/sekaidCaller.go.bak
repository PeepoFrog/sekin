package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	shellquote "github.com/kballard/go-shellquote"
)

// CommandWhitelist maps allowed commands to their executable paths.
var sCallerLog = "/sekai/sCaller.log"

var CommandWhitelist = map[string]string{
	"add-genesis-account":       "/sekaid",
	"collect-gentxs":            "/sekaid",
	"config":                    "/sekaid",
	"debug":                     "/sekaid",
	"export":                    "/sekaid",
	"export-metadata":           "/sekaid",
	"export-minimized-genesis":  "/sekaid",
	"gentx":                     "/sekaid",
	"gentx-claim":               "/sekaid",
	"help":                      "/sekaid",
	"init":                      "/sekaid",
	"keys":                      "/sekaid",
	"new-genesis-from-exported": "/sekaid",
	"query":                     "/sekaid",
	"q":                         "/sekaid",
	"rollback":                  "/sekaid",
	"rosetta":                   "/sekaid",
	"start":                     "/sekaid",
	"status":                    "/sekaid",
	"tendermint":                "/sekaid",
	"testnet":                   "/sekaid",
	"tx":                        "/sekaid",
	"val-address":               "/sekaid",
	"valcons-address":           "/sekaid",
	"validate-genesis":          "/sekaid",
	"version":                   "/sekaid",
}

func main() {
	// Configure logging
	setupLogging()

	log.Println("Server started. Waiting for commands...")

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		commandLine := scanner.Text()
		wg.Add(1)
		go func(command string) {
			defer wg.Done()
			executeCommand(command)
		}(commandLine)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading standard input: %v", err)
	}
}

func setupLogging() {
	logFile, err := os.OpenFile(sCallerLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func executeCommand(commandLine string) {
	// Split command line into command and arguments
	// parts := strings.Fields(commandLine)
	parts, err := shellquote.Split(commandLine)
	if err != nil {
		log.Printf("Error parsing command line: %v", err)

		return
	}

	if len(parts) == 0 {
		log.Printf("Empty command received. Skipping.")

		return
	}

	cmd := parts[0]

	// Validate command against whitelist
	execPath, allowed := CommandWhitelist[cmd]
	if !allowed {
		log.Printf("Command '%s' is not allowed", cmd)

		return
	}

	if cmd == "start" {
		args := append([]string{execPath}, parts...) // Add binary path to args
		if err := syscall.Exec(execPath, args, os.Environ()); err != nil {
			log.Printf("Failed to execute '%s' with syscall.Exec: %v", execPath, err)
		}

		return
	}

	// Execute the command securely
	command := exec.Command(execPath, parts...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Printf("Error executing command '%s': %v", cmd, err)
	} else {
		log.Printf("Command '%s' executed successfully", cmd)
	}
}
