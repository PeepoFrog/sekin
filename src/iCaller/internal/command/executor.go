package command

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
)

func InterxInitCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*InterxInit)
	if !ok {
		return ``, fmt.Errorf("invalid arguments for 'init', error converting args to InterxInit struct\n Args:%v", args)
	}

	cmdMap := make(map[string]interface{})
	cmdMap["addrbook"] = cmdArgs.AddrBook
	cmdMap["cache_dir"] = cmdArgs.CacheDir
	cmdMap["caching_duration"] = cmdArgs.CachingDuration
	cmdMap["download_file_size_limitation"] = cmdArgs.DownloadFileSizeLimitation
	cmdMap["faucet_amounts"] = cmdArgs.FaucetAmounts
	cmdMap["faucet_minimum_amounts"] = cmdArgs.FaucetMinimumAmounts
	cmdMap["faucet_mnemonic"] = cmdArgs.FaucetMnemonic
	cmdMap["faucet_time_limit"] = cmdArgs.FaucetTimeLimit
	cmdMap["fee_amounts"] = cmdArgs.FeeAmounts
	cmdMap["grpc"] = cmdArgs.Grpc
	cmdMap["halted_avg_block_times"] = cmdArgs.HaltedAvgBlockTimes
	cmdMap["home"] = cmdArgs.Home
	cmdMap["max_cache_size"] = cmdArgs.MaxCacheSize
	cmdMap["node_discovery_interx_port"] = cmdArgs.NodeDiscoveryInterxPort
	cmdMap["node_discovery_tendermint_port"] = cmdArgs.NodeDiscoveryTendermintPort
	cmdMap["node_discovery_timeout"] = cmdArgs.NodeDiscoveryTimeout
	cmdMap["node_discovery_use_https"] = cmdArgs.NodeDiscoveryUseHttps
	cmdMap["node_key"] = cmdArgs.NodeKey
	cmdMap["node_type"] = cmdArgs.NodeType
	cmdMap["port"] = cmdArgs.Port
	cmdMap["rpc"] = cmdArgs.Rpc
	cmdMap["seed_node_id"] = cmdArgs.SeedNodeID
	cmdMap["sentry_node_id"] = cmdArgs.SentryNodeID
	cmdMap["serve_https"] = cmdArgs.ServeHttps
	cmdMap["signing_mnemonic"] = cmdArgs.SigningMnemonic
	cmdMap["snapshot_interval"] = cmdArgs.SnapshotInterval
	cmdMap["snapshot_node_id"] = cmdArgs.SnapshotNodeID
	cmdMap["status_sync"] = cmdArgs.StatusSync
	cmdMap["tx_modes"] = cmdArgs.TxModes
	cmdMap["validator_node_id"] = cmdArgs.ValidatorNodeID

	var flagsStr []string = []string{"init"}
	for k, v := range cmdMap {
		if !checkNilInterface(v) && v != "" {
			flagsStr = append(flagsStr, fmt.Sprintf("--%v=%v", k, reflect.Indirect(reflect.ValueOf(v))))
		} else {
			log.Printf("DEBUG: <%v> was not added with <%v> value\n", k, v)
		}
	}

	cmd := exec.Command(ExecPath, flagsStr...)

	log.Printf("DEBUG: formed cmd: %+v", cmd.Args)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func checkNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

func InterxVersionCmd(args interface{}) (string, error) {
	cmd := exec.Command(ExecPath, `version`)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
