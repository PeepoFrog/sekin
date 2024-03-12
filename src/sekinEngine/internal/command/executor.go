package command

// Example
// func SekaiInitCmd(args interface{}) (string, error) {
// 	cmdArgs, ok := args.(*SekaiInit)
// 	re := regexp.MustCompile(`\s+`)

// 	if !ok {
// 		return "", fmt.Errorf("invalid arguments for 'init'")
// 	}

// 	cmd := exec.Command(ExecPath, "init",
// 		"--home", cmdArgs.Home,
// 		"--chain-id", cmdArgs.ChainID,
// 		fmt.Sprintf("%v", re.ReplaceAllString(cmdArgs.Moniker, "_")),
// 		"--log_level", cmdArgs.LogLvl,
// 		"--log_format", cmdArgs.LogFmt,
// 	)

// 	if cmdArgs.Overwrite {
// 		cmd.Args = append(cmd.Args, "--overwrite")
// 	}

// 	log.Printf("DEBUG: SekaiInitCmd: cmd args: %v", cmd.Args)
// 	output, err := cmd.CombinedOutput()
// 	log.Println(string(output))
// 	return string(output), err
// }
