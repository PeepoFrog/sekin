package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	binaryName := "/interxd" // Hardcoded binary name
	args := []string{binaryName}

	// Append "version" if no arguments are provided, else pass all given arguments
	if len(os.Args) == 1 {
		args = append(args, "version")
	} else {
		args = append(args, os.Args[1:]...)
	}
	binary, lookErr := exec.LookPath(binaryName)
	if lookErr != nil {
		panic(lookErr)
	}

	execErr := syscall.Exec(binary, args, os.Environ())
	if execErr != nil {
		panic(execErr)
	}
}
