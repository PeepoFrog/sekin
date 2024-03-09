package command

import (
	"fmt"
	"os/exec"
)

func InterxInitCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(InterxInit)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'init'")
	}
	cmd := exec.Command(ExecPath, "init",
		fmt.Sprintf("-addrbook=%v", cmdArgs.AddrBook),
		fmt.Sprintf("-cache_dir=%v", cmdArgs.CacheDir),
		fmt.Sprintf("-caching_duration=%v", cmdArgs.CachingDuration),
		fmt.Sprintf("-download_file_size_limitation=%v", cmdArgs.DownloadFileSizeLimitation),
		fmt.Sprintf("-faucet_amounts=%v", cmdArgs.FaucetAmounts),
		fmt.Sprintf("-faucet_minimum_amounts=%v", cmdArgs.FaucetMinimumAmounts),
		fmt.Sprintf("-faucet_mnemonic=%v", cmdArgs.FaucetMnemonic),
		fmt.Sprintf("-faucet_time_limit=%v", cmdArgs.FaucetTimeLimit),
		fmt.Sprintf("-fee_amounts=%v", cmdArgs.FeeAmounts),
		fmt.Sprintf("-grpc=%v", cmdArgs.Grpc),
		fmt.Sprintf("-halted_avg_block_times=%v", cmdArgs.HaltedAvgBlockTimes),
		fmt.Sprintf("-home=%v", cmdArgs.Home),
		fmt.Sprintf("-max_cache_size=%v", cmdArgs.MaxCacheSize),
		fmt.Sprintf("-node_discovery_interx_port=%v", cmdArgs.NodeDiscoveryInterxPort),
		fmt.Sprintf("-node_discovery_tendermint_port=%v", cmdArgs.NodeDiscoveryTendermintPort),
		fmt.Sprintf("-node_discovery_timeout=%v", cmdArgs.NodeDiscoveryTimeout),
		fmt.Sprintf("-node_discovery_use_https=%v", cmdArgs.NodeDiscoveryUseHttps),
		fmt.Sprintf("-node_key=%v", cmdArgs.NodeKey),
		fmt.Sprintf("-node_type=%v", cmdArgs.NodeType),
		fmt.Sprintf("-port=%v", cmdArgs.Port),
		fmt.Sprintf("-rpc=%v", cmdArgs.Rpc),
		fmt.Sprintf("-seed_node_id=%v", cmdArgs.SeedNodeID),
		fmt.Sprintf("-sentry_node_id=%v", cmdArgs.SentryNodeID),
		fmt.Sprintf("-serve_https=%v", cmdArgs.ServeHttps),
		fmt.Sprintf("-signing_mnemonic=%v", cmdArgs.SigningMnemonic),
		fmt.Sprintf("-snapshot_interval=%v", cmdArgs.SnapshotInterval),
		fmt.Sprintf("-snapshot_node_id=%v", cmdArgs.SnapshotNodeID),
		fmt.Sprintf("-status_sync=%v", cmdArgs.StatusSync),
		fmt.Sprintf("-tx_modes=%v", cmdArgs.TxModes),
		fmt.Sprintf("-validator_node_id=%v", cmdArgs.ValidatorNodeID),
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func InterxVersionCmd(args interface{}) (string, error) {
	cmd := exec.Command(ExecPath, "version")
	output, err := cmd.CombinedOutput()

	return string(output), err
}
