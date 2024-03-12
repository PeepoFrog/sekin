package command

import (
	"fmt"
	"log"

	cosmosBIP39 "github.com/cosmos/go-bip39"
	kiraMnemonicGen "github.com/kiracore/tools/bip39gen/cmd"
)

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

func InitNewNetwork(args interface{}) (string, error) {
	iArgs, ok := args.(*InitNew)
	if !ok {
		return "", fmt.Errorf("invalid arguments for init")
	}

	var masterMnemonic string
	if iArgs.Mnemonic == "" {
		//generate new mnemonic
		log.Printf("Mnemonic is not specified, generating new one")
		mnemonicToGenerate := kiraMnemonicGen.NewMnemonic()
		mnemonicToGenerate.SetRandomEntropy(24)
		mnemonicToGenerate.Generate()
		masterMnemonic = mnemonicToGenerate.String()
	} else {
		log.Printf("Mnemonic is specified, validating")
		check := cosmosBIP39.IsMnemonicValid(iArgs.Mnemonic)
		if !check {
			return "", fmt.Errorf("mnemonic <%v> is not valid", iArgs.Mnemonic)
		}
		masterMnemonic = iArgs.Mnemonic
	}
	fmt.Printf("Your master mnemonic is: <%v>", masterMnemonic)

	return "Network is initialized", nil
}
