package mnemonicmanager

import (
	"fmt"
	"os"

	vlg "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
	cosmosBIP39 "github.com/cosmos/go-bip39"
)

func GenerateMnemonicsFromMaster(masterMnemonic string) (*vlg.MasterMnemonicSet, error) {
	defaultPrefix := vlg.DefaultPrefix
	defaultPath := vlg.DefaultPath

	mnemonicSet, err := vlg.MasterKeysGen([]byte(masterMnemonic), defaultPrefix, defaultPath, "")
	if err != nil {
		return nil, err
	}

	return &mnemonicSet, nil
}

func ValidateMnemonic(mnemonic string) error {
	check := cosmosBIP39.IsMnemonicValid(mnemonic)
	if !check {
		return fmt.Errorf("mnemonic <%v> is not valid", mnemonic)
	}
	return nil
}

func SetSekaidPrivKeys(mnemonicSet *vlg.MasterMnemonicSet, sekaidHome string) error {
	// TODO path set as variables or constants
	sekaidConfigFolder := sekaidHome + "/config"
	fmt.Println(sekaidConfigFolder)

	//creating sekaid home
	err := os.Mkdir(sekaidHome, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("unable to create <%s> folder, err: %w", sekaidHome, err)
		}
	}
	//creating sekaid's config folder
	err = os.Mkdir(sekaidConfigFolder, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("unable to create <%s> folder, err: %w", sekaidConfigFolder, err)
		}
	}

	err = vlg.GeneratePrivValidatorKeyJson(mnemonicSet.ValidatorValMnemonic, sekaidConfigFolder+"/priv_validator_key.json", vlg.DefaultPrefix, vlg.DefaultPath)
	if err != nil {
		return fmt.Errorf("unable to generate priv_validator_key.json: %w", err)
	}
	err = vlg.GenerateValidatorNodeKeyJson(mnemonicSet.ValidatorNodeMnemonic, sekaidConfigFolder+"/node_key.json", vlg.DefaultPrefix, vlg.DefaultPath)
	if err != nil {
		return fmt.Errorf("unable to generate node_key.json: %w", err)
	}
	return nil
}

// sets empty state of validator into $sekaidHome/data/priv_validator_state.json
func SetEmptyValidatorState(sekaidHome string) error {
	emptyState := `
	{
		"height": "0",
		"round": 0,
		"step": 0
	}`
	sekaidDataFolder := sekaidHome + "/data"
	err := os.Mkdir(sekaidDataFolder, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("unable to create <%s> folder, err: %w", sekaidDataFolder, err)
		}
	}
	// utils.CreateFileWithData(sekaidDataFolder+"/priv_validator_state.json", []byte(emptyState))
	file, err := os.Create(sekaidDataFolder + "/priv_validator_state.json")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(emptyState))
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}
	fmt.Println(emptyState, sekaidDataFolder)
	return nil
}
