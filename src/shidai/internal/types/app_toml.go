package types

type AppConfig struct {
	MinimumGasPrices    string   `toml:"minimum-gas-prices"`
	Pruning             string   `toml:"pruning"`
	PruningKeepRecent   string   `toml:"pruning-keep-recent"`
	PruningInterval     string   `toml:"pruning-interval"`
	HaltHeight          int64    `toml:"halt-height"`
	HaltTime            int64    `toml:"halt-time"`
	MinRetainBlocks     int64    `toml:"min-retain-blocks"`
	InterBlockCache     bool     `toml:"inter-block-cache"`
	IndexEvents         []string `toml:"index-events"`
	IavlCacheSize       int      `toml:"iavl-cache-size"`
	IavlDisableFastNode bool     `toml:"iavl-disable-fastnode"`
	IavlLazyLoading     bool     `toml:"iavl-lazy-loading"`
	AppDBBackend        string   `toml:"app-db-backend"`

	Telemetry TelemetryConfig    `toml:"telemetry"`
	API       APIConfig          `toml:"api"`
	Rosetta   RosettaConfig      `toml:"rosetta"`
	GRPC      GRPCConfig         `toml:"grpc"`
	GRPCWeb   GRPCWebConfig      `toml:"grpc-web"`
	StateSync StateSyncAppConfig `toml:"state-sync"`
	Store     StoreConfig        `toml:"store"`
	Streamers StreamersConfig    `toml:"streamers"`
	Mempool   MempoolAppConfig   `toml:"mempool"`
	Wasm      WasmConfig         `toml:"wasm"`
}

type TelemetryConfig struct {
	ServiceName             string     `toml:"service-name"`
	Enabled                 bool       `toml:"enabled"`
	EnableHostname          bool       `toml:"enable-hostname"`
	EnableHostnameLabel     bool       `toml:"enable-hostname-label"`
	EnableServiceLabel      bool       `toml:"enable-service-label"`
	PrometheusRetentionTime int        `toml:"prometheus-retention-time"`
	GlobalLabels            [][]string `toml:"global-labels"`
}

type APIConfig struct {
	Enable             bool   `toml:"enable"`
	Swagger            bool   `toml:"swagger"`
	Address            string `toml:"address"`
	MaxOpenConnections int    `toml:"max-open-connections"`
	RPCReadTimeout     int    `toml:"rpc-read-timeout"`
	RPCWriteTimeout    int    `toml:"rpc-write-timeout"`
	RPCMaxBodyBytes    int    `toml:"rpc-max-body-bytes"`
	EnabledUnsafeCORS  bool   `toml:"enabled-unsafe-cors"`
}

type RosettaConfig struct {
	Enable              bool   `toml:"enable"`
	Address             string `toml:"address"`
	Blockchain          string `toml:"blockchain"`
	Network             string `toml:"network"`
	Retries             int    `toml:"retries"`
	Offline             bool   `toml:"offline"`
	EnableFeeSuggestion bool   `toml:"enable-fee-suggestion"`
	GasToSuggest        int    `toml:"gas-to-suggest"`
	DenomToSuggest      string `toml:"denom-to-suggest"`
}

type GRPCConfig struct {
	Enable         bool   `toml:"enable"`
	Address        string `toml:"address"`
	MaxRecvMsgSize string `toml:"max-recv-msg-size"`
	MaxSendMsgSize string `toml:"max-send-msg-size"`
}

type GRPCWebConfig struct {
	Enable           bool   `toml:"enable"`
	Address          string `toml:"address"`
	EnableUnsafeCORS bool   `toml:"enable-unsafe-cors"`
}

type StateSyncAppConfig struct {
	SnapshotInterval   int `toml:"snapshot-interval"`
	SnapshotKeepRecent int `toml:"snapshot-keep-recent"`
}

type StoreConfig struct {
	Streamers []string `toml:"streamers"`
}

type StreamersConfig struct {
	File StreamersFileConfig `toml:"file"`
}

type StreamersFileConfig struct {
	Keys            []string `toml:"keys"`
	WriteDir        string   `toml:"write_dir"`
	Prefix          string   `toml:"prefix"`
	OutputMetadata  string   `toml:"output-metadata"`
	StopNodeOnError string   `toml:"stop-node-on-error"`
	Fsync           string   `toml:"fsync"`
}

type MempoolAppConfig struct {
	MaxTxs int `toml:"max-txs"`
}

type WasmConfig struct {
	QueryGasLimit int `toml:"query_gas_limit"`
	LRUSize       int `toml:"lru_size"`
}

func NewDefaultAppConfig() *AppConfig {
	return &AppConfig{
		// general fields
		MinimumGasPrices:  "0stake",
		Pruning:           "custom",
		PruningKeepRecent: "2",
		// PruningKeepEvery:    "100",
		PruningInterval:     "10",
		HaltHeight:          0,
		HaltTime:            0,
		MinRetainBlocks:     0,
		InterBlockCache:     true,
		IndexEvents:         []string{},
		IavlCacheSize:       781250,
		IavlDisableFastNode: true,
		AppDBBackend:        "",
		// [telemetry]
		Telemetry: TelemetryConfig{
			ServiceName:             "",
			Enabled:                 false,
			EnableHostname:          false,
			EnableHostnameLabel:     false,
			EnableServiceLabel:      false,
			PrometheusRetentionTime: 0,
			GlobalLabels:            [][]string{},
		},
		// [api]
		API: APIConfig{
			Enable:             true, // api should be true for new interx
			Swagger:            false,
			Address:            "tcp://0.0.0.0:1317", // 0.0.0.0 or localhost?
			MaxOpenConnections: 1000,
			RPCReadTimeout:     10,
			RPCWriteTimeout:    0,
			RPCMaxBodyBytes:    1000000,
			EnabledUnsafeCORS:  false,
		},
		// [rosetta]
		Rosetta: RosettaConfig{
			Enable:              false,
			Address:             ":8080",
			Blockchain:          "app",
			Network:             "network",
			Retries:             3,
			Offline:             false,
			EnableFeeSuggestion: false,
			GasToSuggest:        200000,
			DenomToSuggest:      "ukex", //default uatom
		},
		// [grpc]
		GRPC: GRPCConfig{
			Enable:         true,
			Address:        "0.0.0.0:9090", // 0.0.0.0 or localhost?
			MaxRecvMsgSize: "10485760",     // default value
			MaxSendMsgSize: "2147483647",   // default value
		},
		// [grpc-web]
		GRPCWeb: GRPCWebConfig{
			Enable:           true,
			Address:          "0.0.0.0:9091", // 0.0.0.0 or localhost?
			EnableUnsafeCORS: false,
		},
		// [state-sync]
		StateSync: StateSyncAppConfig{
			SnapshotInterval:   200,
			SnapshotKeepRecent: 2,
		},
		// [store]
		Store: StoreConfig{
			Streamers: []string{},
		},
		// [streamers]
		Streamers: StreamersConfig{
			// [streamers.file]
			StreamersFileConfig{
				Keys:            []string{"*"},
				WriteDir:        "",
				Prefix:          "",
				OutputMetadata:  "true",
				StopNodeOnError: "true",
				Fsync:           "false",
			},
		},
		// [mempool]
		Mempool: MempoolAppConfig{
			MaxTxs: 5000, // default value
		},
		// [wasm]
		Wasm: WasmConfig{
			QueryGasLimit: 300000,
			LRUSize:       0,
		},
	}

}
