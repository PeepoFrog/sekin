package sekaihandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	mnemonicsgenerator "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
	httpexecutor "github.com/kiracore/sekin/src/shidai/internal/http_executor"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
	mnemonicmanager "github.com/kiracore/sekin/src/shidai/internal/mnemonic_manager"

	configconstructor "github.com/kiracore/sekin/src/shidai/internal/sekai_handler/config_constructor"
	genesishandler "github.com/kiracore/sekin/src/shidai/internal/sekai_handler/genesis_handler"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"github.com/kiracore/sekin/src/shidai/internal/utils"
)

var log = logger.GetLogger()

func InitSekaiJoiner(ctx context.Context, tc *configconstructor.TargetSeedKiraConfig, masterMnemonicSet *mnemonicsgenerator.MasterMnemonicSet) error {
	cmd := httpexecutor.CommandRequest{
		Command: "init",
		Args: map[string]interface{}{
			"home":     types.SEKAI_HOME,
			"chain-id": "initnet-1",
			"moniker":  "validator node",
		},
	}
	_, err := httpexecutor.ExecuteCallerCommand("sekai.local", "8080", "POST", cmd)
	if err != nil {
		return fmt.Errorf("unable execute <%v> request, error: %w", cmd, err)
	}
	err = setSekaidKeys(masterMnemonicSet)
	if err != nil {
		return fmt.Errorf("unable to set sekai keys: %w", err)
	}
	genesis, err := genesishandler.GetVerifiedGenesisFile(ctx, tc.IpAddress, tc.SekaidRPCPort, tc.InterxPort)
	if err != nil {
		return fmt.Errorf("unable to receive genesis file: %w", err)
	}
	err = os.WriteFile(fmt.Sprintf("%v/config/genesis.json", types.SEKAI_HOME), genesis, 0644)
	if err != nil {
		return fmt.Errorf("cant write genesis.json file: %w", err)
	}

	err = configconstructor.FormSekaiJoinerConfigs(tc)
	if err != nil {
		return fmt.Errorf("unable retrieve join information from <%s>, error: %w", tc.IpAddress, err)
	}

	return nil
}

func setSekaidKeys(masterMnemonicSet *mnemonicsgenerator.MasterMnemonicSet) error {
	err := mnemonicmanager.SetSekaidPrivKeys(masterMnemonicSet, types.SEKAI_HOME)
	if err != nil {
		return fmt.Errorf("unable to set sekaid keys: %w", err)
	}
	err = mnemonicmanager.SetEmptyValidatorState(types.SEKAI_HOME)
	if err != nil {
		return fmt.Errorf("unable to set empty validator state : %w", err)
	}

	_, err = utils.AddKeyToKeyring("validator", string(masterMnemonicSet.ValidatorAddrMnemonic), types.SEKAI_HOME, "test")
	if err != nil {
		return fmt.Errorf("unable to add validator key to keyring: %w", err)
	}
	return nil
}

func StartSekai() error {
	cmd := httpexecutor.CommandRequest{
		Command: "start",
		Args: map[string]interface{}{
			"home": types.SEKAI_HOME,
		},
	}
	_, err := httpexecutor.ExecuteCallerCommand(types.SEKAI_CONTAINER_ADDRESS, "8080", "POST", cmd)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Debug("sekai started")
		} else {
			return fmt.Errorf("unable execute <%v> request, error: %w", cmd, err)
		}
	}

	return nil
}
