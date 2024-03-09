package command

const ExecPath = "/interx"

type CommandHandler func(interface{}) (string, error)

var CommandMapping = map[string]struct {
	ArgsStruct interface{}
	Handler    CommandHandler
}{
	"init":    {ArgsStruct: InterxInit{}, Handler: InterxInitCmd},
	"version": {ArgsStruct: InterxVersion{}, Handler: InterxVersionCmd},
}

type InterxVersion struct {
}

type InterxInit struct {
	AddrBook                    string `json:"addrbook,omitempty"`
	CacheDir                    string `json:"cache_dir,omitempty"`
	CachingDuration             int    `json:"caching_duration,omitempty"`
	DownloadFileSizeLimitation  string `json:"download_file_size_limitation,omitempty"`
	FaucetAmounts               string `json:"faucet_amounts,omitempty"`
	FaucetMinimumAmounts        string `json:"faucet_minimum_amounts,omitempty"`
	FaucetMnemonic              string `json:"faucet_mnemonic,omitempty"`
	FaucetTimeLimit             int    `json:"faucet_time_limit,omitempty"`
	FeeAmounts                  string `json:"fee_amounts,omitempty"`
	Grpc                        string `json:"grpc"`
	HaltedAvgBlockTimes         int    `json:"halted_avg_block_times,omitempty"`
	Home                        string `json:"home,omitempty"`
	MaxCacheSize                string `json:"max_cache_size,omitempty"`
	NodeDiscoveryInterxPort     string `json:"node_discovery_interx_port,omitempty"`
	NodeDiscoveryTendermintPort string `json:"node_discovery_tendermint_port,omitempty"`
	NodeDiscoveryTimeout        string `json:"node_discovery_timeout,omitempty"`
	NodeDiscoveryUseHttps       bool   `json:"node_discovery_use_https,omitempty"`
	NodeKey                     string `json:"node_key,omitempty"`
	NodeType                    string `json:"node_type,omitempty"`
	Port                        string `json:"port"`
	Rpc                         string `json:"rpc"`
	SeedNodeID                  string `json:"seed_node_id,omitempty"`
	SentryNodeID                string `json:"sentry_node_id,omitempty"`
	ServeHttps                  bool   `json:"serve_https,omitempty"`
	SigningMnemonic             string `json:"signing_mnemonic,omitempty"`
	SnapshotInterval            uint   `json:"snapshot_interval,omitempty"`
	SnapshotNodeID              string `json:"snapshot_node_id,omitempty"`
	StatusSync                  int    `json:"status_sync,omitempty"`
	TxModes                     string `json:"tx_modes,omitempty"`
	ValidatorNodeID             string `json:"validator_node_id,omitempty"`
}
