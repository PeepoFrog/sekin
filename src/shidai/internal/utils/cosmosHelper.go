package utils

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"go.uber.org/zap"
)

func AddKeyToKeyring(keyName, mnemonic, homeFolder, keyringType string) (*keyring.Record, error) {
	if keyName == "" {
		return nil, fmt.Errorf("key name cannot be empty")
	}
	check := bip39.IsMnemonicValid(mnemonic)
	if !check {
		return nil, fmt.Errorf("mnemonic is not valid <%v>", mnemonic)
	}
	log.Debug("received mnemonic is valid: ", zap.String("mnemonic", mnemonic))

	if keyringType == "" {
		keyringType = keyring.BackendOS // setting up default value for keyring = "os" check AddKeyringFlags() from sekai
	}

	registry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(registry)
	kb, err := keyring.New(
		sdk.KeyringServiceName(), // Keyring name
		keyringType,              // Backend type
		homeFolder,               // Keys directory path
		os.Stdin,                 // io.Reader for entropy
		marshaler,                // codec.Codec for encoding/decoding
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new keyring: %w", err)
	}
	// from cosmosSdk cmd AddKeyCommand() from cosmosSDK/client/keys/add.go
	// default values from cosmosSDK (same for sekai)
	coinType := sdk.GetConfig().GetCoinType()
	var account uint32 = 0
	var index uint32 = 0
	algoStr := string(hd.Secp256k1Type)
	keyringAlgos, _ := kb.SupportedAlgorithms()
	log.Debug(fmt.Sprintf("default values for algo string: %v, %v, %v, %v, %v", coinType, account, index, algoStr, keyringAlgos))

	algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
	if err != nil {
		return nil, fmt.Errorf("error creating new signing algorithm: %w", err)
	}

	hdPath := hd.CreateHDPath(coinType, account, index).String()
	log.Debug(fmt.Sprintf("hdPath: %v", hdPath))

	k, err := kb.NewAccount(keyName, mnemonic, "", hdPath, algo)
	if err != nil {
		return nil, fmt.Errorf("error creating new account: %w", err)
	}
	k.GetAddress()
	return k, nil
}
