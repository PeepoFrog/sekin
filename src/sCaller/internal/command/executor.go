package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func SekaiInitCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiInit)

	if !ok {
		return "", fmt.Errorf("invalid arguments for 'init'")
	}

	cmd := exec.Command(ExecPath, "init",
		"--home", cmdArgs.Home,
		"--chain-id", cmdArgs.ChainID,
		fmt.Sprintf("%q", cmdArgs.Moniker),
		"--log_level", cmdArgs.LogLvl,
		"--log_format", cmdArgs.LogFmt,
	)

	if cmdArgs.Overwrite {
		cmd.Args = append(cmd.Args, "--overwrite")
	}

	log.Printf("DEBUG: SekaiInitCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiVersionCmd(interface{}) (string, error) {
	cmd := exec.Command(ExecPath, "version")
	log.Printf("DEBUG: SekaiVersionCmd: cmd: %v", cmd)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func SekaidKeysAddCmd(args interface{}) (string, error) {
	log.Printf("DEBUG: SekaidKeysAddCmd: in args: %v", args)
	cmdArgs, ok := args.(*SekaidKeysAdd)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'keys-add'")
	}

	cmd := exec.Command(ExecPath, "keys", "add", cmdArgs.Address,
		"--keyring-backend", cmdArgs.Keyring,
		"--home", cmdArgs.Home,
		"--log_format", cmdArgs.LogFmt,
		"--log_level", cmdArgs.LogLvl,
	)

	if cmdArgs.Output != "" {
		cmd.Args = append(cmd.Args, "--output", cmdArgs.Output)
	}
	if cmdArgs.Recover {
		cmd.Args = append(cmd.Args, "--recover")
	}
	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}

	log.Printf("DEBUG: SekaidKeysAddCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiAddGenesisAccCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiAddGenesisAcc)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'add-genesis-account'")
	}

	cmd := exec.Command(ExecPath, "add-genesis-account", cmdArgs.Address, strings.Join(cmdArgs.Coins, ","), "--home", cmdArgs.Home, "--keyring-backend", cmdArgs.Keyring, "--log_format", cmdArgs.LogFmt, "--log_level", cmdArgs.LogLvl)
	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}

	log.Printf("DEBUG: SekaiAddGenesisAccCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return string(output), err
}

func SekaiGentxClaimCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaiGentxClaim)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'gentx-claim'")
	}
	cmd := exec.Command(
		ExecPath, "gentx-claim", cmdArgs.Address,
		"--keyring-backend", cmdArgs.Keyring,
		"--moniker", fmt.Sprintf("%q", cmdArgs.Moniker),
		"--pubkey", cmdArgs.PubKey,
		"--home", cmdArgs.Home,
		"--log_format", cmdArgs.LogFmt,
		"--log_level", cmdArgs.LogLvl)

	if cmdArgs.Trace {
		cmd.Args = append(cmd.Args, "--trace")
	}
	log.Printf("DEBUG: SekaiGentxClaimCmd: cmd args: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))

	return string(output), err
}

func SekaidStartCmd(args interface{}) (string, error) {
	cmdArgs, ok := args.(*SekaidStart)
	if !ok {
		return "", fmt.Errorf("invalid arguments for 'start'")
	}

	argv := []string{"sekaid", "start", "--home", cmdArgs.Home, "--rpc.laddr", cmdArgs.RPC.LAddr}
	if cmdArgs.Trace {
		argv = append(argv, "--trace")
	}
	if cmdArgs.Moniker != "" {
		argv = append(argv, "--moniker", cmdArgs.Moniker)
	}
	if cmdArgs.Consensus.CreateEmptyBlocks {
		argv = append(argv, "--consensus.create_empty_blocks=true")
	} else {
		argv = append(argv, "--consensus.create_empty_blocks=false")
	}
	if cmdArgs.Consensus.CreateEmptyBlocksInterval != "" {
		argv = append(argv, "--consensus.create_empty_blocks_interval", cmdArgs.Consensus.CreateEmptyBlocksInterval)
	}
	if cmdArgs.Consensus.DoubleSignCheckHeight > 0 {
		argv = append(argv, "--consensus.double_sign_check_height", fmt.Sprintf("%d", cmdArgs.Consensus.DoubleSignCheckHeight))
	}
	if cmdArgs.CPUProfile != "" {
		argv = append(argv, "--cpu-profile", cmdArgs.CPUProfile)
	}
	if cmdArgs.Database.Backend != "" {
		argv = append(argv, "--db_backend", cmdArgs.Database.Backend)
	}
	if cmdArgs.Database.Dir != "" {
		argv = append(argv, "--db_dir", cmdArgs.Database.Dir)
	}
	if !cmdArgs.FastSync {
		argv = append(argv, "--fast_sync=false")
	}
	if cmdArgs.GenesisHash != "" {
		argv = append(argv, "--genesis_hash", cmdArgs.GenesisHash)
	}
	if cmdArgs.GRPCOnly {
		argv = append(argv, "--grpc-only")
	}
	if cmdArgs.GRPCWebAddress != "" {
		argv = append(argv, "--grpc-web.address", cmdArgs.GRPCWebAddress)
	}
	if !cmdArgs.GRPCWebEnable {
		argv = append(argv, "--grpc-web.enable=false")
	}
	if cmdArgs.GRPCAddress != "" {
		argv = append(argv, "--grpc.address", cmdArgs.GRPCAddress)
	}
	if !cmdArgs.GRPCEnable {
		argv = append(argv, "--grpc.enable=false")
	}
	if cmdArgs.HaltHeight > 0 {
		argv = append(argv, "--halt-height", fmt.Sprintf("%d", cmdArgs.HaltHeight))
	}
	if cmdArgs.HaltTime > 0 {
		argv = append(argv, "--halt-time", fmt.Sprintf("%d", cmdArgs.HaltTime))
	}
	if cmdArgs.IAVLDisableFastNode {
		argv = append(argv, "--iavl-disable-fastnode")
	}
	if !cmdArgs.InterBlockCache {
		argv = append(argv, "--inter-block-cache=false")
	}
	if cmdArgs.InvCheckPeriod > 0 {
		argv = append(argv, "--inv-check-period", fmt.Sprintf("%d", cmdArgs.InvCheckPeriod))
	}
	if cmdArgs.MinRetainBlocks > 0 {
		argv = append(argv, "--min-retain-blocks", fmt.Sprintf("%d", cmdArgs.MinRetainBlocks))
	}
	if cmdArgs.MinimumGasPrices != "" {
		argv = append(argv, "--minimum-gas-prices", cmdArgs.MinimumGasPrices)
	}
	if cmdArgs.P2PExternalAddress != "" {
		argv = append(argv, "--p2p.external-address", cmdArgs.P2PExternalAddress)
	}
	if cmdArgs.P2PLAddr != "" {
		argv = append(argv, "--p2p.laddr", cmdArgs.P2PLAddr)
	}
	if cmdArgs.P2PPersistentPeers != "" {
		argv = append(argv, "--p2p.persistent_peers", cmdArgs.P2PPersistentPeers)
	}
	if !cmdArgs.P2PPEX {
		argv = append(argv, "--p2p.pex=false")
	}
	if cmdArgs.P2PPrivatePeerIDs != "" {
		argv = append(argv, "--p2p.private_peer_ids", cmdArgs.P2PPrivatePeerIDs)
	}
	if cmdArgs.P2PSeedMode {
		argv = append(argv, "--p2p.seed_mode")
	}
	if cmdArgs.P2PSeeds != "" {
		argv = append(argv, "--p2p.seeds", cmdArgs.P2PSeeds)
	}
	if cmdArgs.P2PUnconditionalPeerIDs != "" {
		argv = append(argv, "--p2p.unconditional_peer_ids", cmdArgs.P2PUnconditionalPeerIDs)
	}
	if cmdArgs.P2PUPnP {

		argv = append(argv, "--p2p.upnp")
	}
	if cmdArgs.PrivValidatorLAddr != "" {
		argv = append(argv, "--priv_validator_laddr", cmdArgs.PrivValidatorLAddr)
	}
	if cmdArgs.ProxyApp != "" {
		argv = append(argv, "--proxy_app", cmdArgs.ProxyApp)
	}
	if cmdArgs.Pruning != "" {
		argv = append(argv, "--pruning", cmdArgs.Pruning)
	}
	if cmdArgs.PruningInterval > 0 {
		argv = append(argv, "--pruning-interval", fmt.Sprintf("%d", cmdArgs.PruningInterval))
	}
	if cmdArgs.PruningKeepEvery > 0 {
		argv = append(argv, "--pruning-keep-every", fmt.Sprintf("%d", cmdArgs.PruningKeepEvery))
	}
	if cmdArgs.PruningKeepRecent > 0 {
		argv = append(argv, "--pruning-keep-recent", fmt.Sprintf("%d", cmdArgs.PruningKeepRecent))
	}
	if cmdArgs.RPCGRPCAddr != "" {
		argv = append(argv, "--rpc.grpc_laddr", cmdArgs.RPCGRPCAddr)
	}
	if cmdArgs.RPCPprofLAddr != "" {
		argv = append(argv, "--rpc.pprof_laddr", cmdArgs.RPCPprofLAddr)
	}
	if cmdArgs.RPCUnsafe {
		argv = append(argv, "--rpc.unsafe")
	}
	if cmdArgs.StateSyncSnapshotInterval > 0 {
		argv = append(argv, "--state-sync.snapshot-interval", fmt.Sprintf("%d", cmdArgs.StateSyncSnapshotInterval))
	}
	if cmdArgs.StateSyncSnapshotKeepRecent > 0 {
		argv = append(argv, "--state-sync.snapshot-keep-recent", fmt.Sprintf("%d", cmdArgs.StateSyncSnapshotKeepRecent))
	}
	if cmdArgs.TraceStore != "" {
		argv = append(argv, "--trace-store", cmdArgs.TraceStore)
	}
	if cmdArgs.Transport != "" {
		argv = append(argv, "--transport", cmdArgs.Transport)
	}
	for _, skip := range cmdArgs.UnsafeSkipUpgrades {
		argv = append(argv, "--unsafe-skip-upgrades", fmt.Sprintf("%d", skip))
	}
	if cmdArgs.WithTendermint {
		argv = append(argv, "--with-tendermint")
	}
	if cmdArgs.XCrisisSkipAssertInvariants {
		argv = append(argv, "--x-crisis-skip-assert-invariants")
	}

	env := os.Environ()
	err := syscall.Exec(ExecPath, argv, env)
	return "", err

}
