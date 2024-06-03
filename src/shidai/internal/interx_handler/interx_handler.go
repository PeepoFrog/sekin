package interxhandler

import (
	"context"
	"errors"
	"fmt"
	"io"

	mnemonicsgenerator "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"

	httpexecutor "github.com/kiracore/sekin/src/shidai/internal/http_executor"
	"github.com/kiracore/sekin/src/shidai/internal/logger"
	"github.com/kiracore/sekin/src/shidai/internal/types"
)

var log = logger.GetLogger()

func InitInterx(ctx context.Context, masterMnemonicSet *mnemonicsgenerator.MasterMnemonicSet) error {
	signerMnemonic := string(masterMnemonicSet.SignerAddrMnemonic)
	nodeType := "validator"

	grpcSekaidAddress := fmt.Sprintf("dns:///%v:%v", types.SEKAI_CONTAINER_ADDRESS, types.DEFAULT_GRPC_PORT)
	rpcSekaidAddress := fmt.Sprintf("http://%v:%v", types.SEKAI_CONTAINER_ADDRESS, types.DEFAULT_RPC_PORT)
	cmd := httpexecutor.CommandRequest{
		Command: "init",
		Args: map[string]interface{}{
			"home":              types.INTERX_HOME,
			"grpc":              grpcSekaidAddress,
			"rpc":               rpcSekaidAddress,
			"node_type":         nodeType,
			"faucet_mnemonic":   signerMnemonic,
			"signing_mnemonic":  signerMnemonic,
			"port":              types.DEFAULT_INTERX_PORT,
			"validator_node_id": string(masterMnemonicSet.ValidatorNodeId),
		},
	}
	out, err := httpexecutor.ExecuteCallerCommand(types.INTERX_CONTAINER_ADDRESS, "8081", "POST", cmd)
	if err != nil {
		return fmt.Errorf("unable execute <%v> request, error: %w", cmd, err)
	}
	log.Info(string(out))

	return nil
}

func StartInterx() error {
	cmd := httpexecutor.CommandRequest{
		Command: "start",
		Args: map[string]interface{}{
			"home": types.INTERX_HOME,
		},
	}
	_, err := httpexecutor.ExecuteCallerCommand(types.INTERX_CONTAINER_ADDRESS, "8081", "POST", cmd)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Debug("interx started")
		} else {
			return fmt.Errorf("unable execute <%v> request, error: %w", cmd, err)
		}
	}

	return nil
}
