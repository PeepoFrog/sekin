package command

type (
	Handler func(interface{}) (string, error)

	SekaiInit struct {
		ChainID   string `json:"chain-id"`
		Moniker   string `json:"moniker"`
		Home      string `json:"home"`
		LogFmt    string `json:"log_format"`
		LogLvl    string `json:"log_level"`
		Trace     bool   `json:"trace"`
		Overwrite bool   `json:"overwrite"`
	}

	SekaidKeysAdd struct {
		Address string `json:"address"`
		Keyring string `json:"keyring-backend"`
		Home    string `json:"home"`
		LogFmt  string `json:"log_format"`
		LogLvl  string `json:"log_level"`
		Output  string `json:"output"`
		Seed    string `json:"seed"`
		Trace   bool   `json:"trace"`
		Recover bool   `json:"recover"`
	}

	SekaiVersion struct {
	}

	SekaiAddGenesisAcc struct {
		Address string   `json:"address"`
		Home    string   `json:"home"`
		Keyring string   `json:"keyring-backend"`
		LogFmt  string   `json:"log_format"`
		LogLvl  string   `json:"log_level"`
		Trace   bool     `json:"trace"`
		Coins   []string `json:"coins"`
	}

	SekaiGentxClaim struct {
		Address string `json:"address"`
		Keyring string `json:"keyring-backend"`
		Moniker string `json:"moniker"`
		PubKey  string `json:"pubkey"`
		Home    string `json:"home"`
		LogFmt  string `json:"log_format"`
		LogLvl  string `json:"log_level"`
		Trace   bool   `json:"trace"`
	}

	SekaidStart struct {
		Moniker                     string          `json:"moniker"` //+
		Home                        string          `json:"home"`    //+
		Trace                       bool            `json:"trace"`   //+
		LogFormat                   string          `json:"log_format"`
		LogLevel                    string          `json:"log_level"`
		Consensus                   ConsensusConfig `json:"consensus"`
		GRPC                        GRPCConfig      `json:"grpc"`
		P2P                         P2PConfig       `json:"p2p"`
		Pruning                     PruningConfig   `json:"pruning"`
		RPC                         RPCConfig       `json:"rpc"`
		StateSync                   StateSyncConfig `json:"state_sync"`
		Database                    DatabaseConfig  `json:"database"`
		ABCI                        ABCIConfig      `json:"abci"`
		CPUProfile                  string          `json:"cpu_profile"`
		GenesisHash                 string          `json:"genesis_hash"`
		MinRetainBlocks             uint            `json:"min_retain_blocks"`
		MinimumGasPrices            string          `json:"minimum_gas_prices"`
		UnsafeSkipUpgrades          []int           `json:"unsafe_skip_upgrades"`
		PrivValidatorLAddr          string          `json:"priv_validator_laddr"`
		ProxyApp                    string          `json:"proxy_app"`
		Transport                   string          `json:"transport"`
		XCrisisSkipAssertInvariants bool            `json:"x_crisis_skip_assert_invariants"`
	}

	ABCIConfig struct {
		ABCIChainID string `json:"abci_chain_id"`
		ABCIAddress string `json:"abci_address"`
	}

	ConsensusConfig struct { //+
		CreateEmptyBlocks         bool   `json:"create_empty_blocks"`          //+
		CreateEmptyBlocksInterval string `json:"create_empty_blocks_interval"` //+
		DoubleSignCheckHeight     int    `json:"double_sign_check_height"`     //+
	}

	GRPCConfig struct {
		Only       bool   `json:"only"`
		WebAddress string `json:"web_address"`
		WebEnable  bool   `json:"web_enable"`
		Address    string `json:"address"`
		Enable     bool   `json:"enable"`
	}

	P2PConfig struct {
		ExternalAddress      string `json:"external_address"`
		LAddr                string `json:"p2p_laddr"`
		PersistentPeers      string `json:"persistent_peers"`
		PEX                  bool   `json:"pex"`
		PrivatePeerIDs       string `json:"private_peer_ids"`
		SeedMode             bool   `json:"seed_mode"`
		Seeds                string `json:"seeds"`
		UnconditionalPeerIDs string `json:"unconditional_peer_ids"`
		UPnP                 bool   `json:"upnp"`
	}

	PruningConfig struct {
		Strategy   string `json:"pruning_strategy"`
		Interval   uint   `json:"pruning_interval"`
		KeepEvery  uint   `json:"pruning_keep_every"`
		KeepRecent uint   `json:"pruning_keep_recent"`
	}

	RPCConfig struct {
		GRPCLAddr  string `json:"grpc_laddr"`
		LAddr      string `json:"rpc_laddr"` //+
		PprofLAddr string `json:"pprof_laddr"`
		Unsafe     bool   `json:"unsafe"`
	}

	StateSyncConfig struct {
		SnapshotInterval   uint   `json:"snapshot_interval"`
		SnapshotKeepRecent uint32 `json:"snapshot_keep_recent"`
	}

	DatabaseConfig struct {
		Backend string `json:"db_backend"`
		Dir     string `json:"db_dir"`
	}
)

const ExecPath = "/sekaid"

var CommandMapping = map[string]struct {
	ArgsStruct func() interface{}
	Handler    Handler
}{
	"init":                {ArgsStruct: func() interface{} { return &SekaiInit{} }, Handler: SekaiInitCmd},
	"version":             {ArgsStruct: func() interface{} { return &SekaiVersion{} }, Handler: SekaiVersionCmd},
	"add-genesis-account": {ArgsStruct: func() interface{} { return &SekaiAddGenesisAcc{} }, Handler: SekaiAddGenesisAccCmd},
	"gentx-claim":         {ArgsStruct: func() interface{} { return &SekaiGentxClaim{} }, Handler: SekaiGentxClaimCmd},
	"keys-add":            {ArgsStruct: func() interface{} { return &SekaidKeysAdd{} }, Handler: SekaidKeysAddCmd},
	"start":               {ArgsStruct: func() interface{} { return &SekaidStart{} }, Handler: SekaidStartCmd},
}
