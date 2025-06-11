package types

type Config struct {
	ProxyApp               string                `toml:"proxy_app"`
	Moniker                string                `toml:"moniker"`
	BlockSync              bool                  `toml:"block_sync"`
	DBBackend              string                `toml:"db_backend"`
	DBDir                  string                `toml:"db_dir"`
	LogLevel               string                `toml:"log_level"`
	LogFormat              string                `toml:"log_format"`
	GenesisFile            string                `toml:"genesis_file"`
	PrivValidatorKeyFile   string                `toml:"priv_validator_key_file"`
	PrivValidatorStateFile string                `toml:"priv_validator_state_file"`
	PrivValidatorLaddr     string                `toml:"priv_validator_laddr"`
	NodeKeyFile            string                `toml:"node_key_file"`
	ABCI                   string                `toml:"abci"`
	FilterPeers            bool                  `toml:"filter_peers"`
	RPC                    RPCConfig             `toml:"rpc"`
	P2P                    P2PConfig             `toml:"p2p"`
	Mempool                MempoolConfig         `toml:"mempool"`
	StateSync              StateSyncConfig       `toml:"statesync"`
	BlockSyncSection       BlockSyncConfig       `toml:"blocksync"`
	Consensus              ConsensusConfig       `toml:"consensus"`
	Storage                StorageConfig         `toml:"storage"`
	TxIndex                TxIndexConfig         `toml:"tx_index"`
	Instrumentation        InstrumentationConfig `toml:"instrumentation"`
}

type RPCConfig struct {
	Laddr                                string   `toml:"laddr"`
	CORSAllowedOrigins                   []string `toml:"cors_allowed_origins"`
	CORSAllowedMethods                   []string `toml:"cors_allowed_methods"`
	CORSAllowedHeaders                   []string `toml:"cors_allowed_headers"`
	GRPCLaddr                            string   `toml:"grpc_laddr"`
	GRPCMaxOpenConnections               int      `toml:"grpc_max_open_connections"`
	Unsafe                               bool     `toml:"unsafe"`
	MaxOpenConnections                   int      `toml:"max_open_connections"`
	MaxSubscriptionClients               int      `toml:"max_subscription_clients"`
	MaxSubscriptionsPerClient            int      `toml:"max_subscriptions_per_client"`
	ExperimentalSubscriptionBufferSize   int      `toml:"experimental_subscription_buffer_size"`
	ExperimentalWebsocketWriteBufferSize int      `toml:"experimental_websocket_write_buffer_size"`
	ExperimentalCloseOnSlowClient        bool     `toml:"experimental_close_on_slow_client"`
	TimeoutBroadcastTxCommit             string   `toml:"timeout_broadcast_tx_commit"`
	MaxBodyBytes                         int      `toml:"max_body_bytes"`
	MaxHeaderBytes                       int      `toml:"max_header_bytes"`
	TLSCertFile                          string   `toml:"tls_cert_file"`
	TLSKeyFile                           string   `toml:"tls_key_file"`
	PprofLaddr                           string   `toml:"pprof_laddr"`
}

type P2PConfig struct {
	Laddr                        string `toml:"laddr"`
	ExternalAddress              string `toml:"external_address"`
	Seeds                        string `toml:"seeds"`
	PersistentPeers              string `toml:"persistent_peers"`
	UPNP                         bool   `toml:"upnp"`
	AddrBookFile                 string `toml:"addr_book_file"`
	AddrBookStrict               bool   `toml:"addr_book_strict"`
	MaxNumInboundPeers           int    `toml:"max_num_inbound_peers"`
	MaxNumOutboundPeers          int    `toml:"max_num_outbound_peers"`
	UnconditionalPeerIDs         string `toml:"unconditional_peer_ids"`
	PersistentPeersMaxDialPeriod string `toml:"persistent_peers_max_dial_period"`
	FlushThrottleTimeout         string `toml:"flush_throttle_timeout"`
	MaxPacketMsgPayloadSize      int    `toml:"max_packet_msg_payload_size"`
	SendRate                     int64  `toml:"send_rate"`
	RecvRate                     int64  `toml:"recv_rate"`
	PEX                          bool   `toml:"pex"`
	SeedMode                     bool   `toml:"seed_mode"`
	PrivatePeerIDs               string `toml:"private_peer_ids"`
	AllowDuplicateIP             bool   `toml:"allow_duplicate_ip"`
	HandshakeTimeout             string `toml:"handshake_timeout"`
	DialTimeout                  string `toml:"dial_timeout"`
}

type MempoolConfig struct {
	Version               string `toml:"version"`
	Recheck               bool   `toml:"recheck"`
	Broadcast             bool   `toml:"broadcast"`
	WALDir                string `toml:"wal_dir"`
	Size                  int    `toml:"size"`
	MaxTxsBytes           int64  `toml:"max_txs_bytes"`
	CacheSize             int    `toml:"cache_size"`
	KeepInvalidTxsInCache bool   `toml:"keep-invalid-txs-in-cache"`
	MaxTxBytes            int    `toml:"max_tx_bytes"`
	MaxBatchBytes         int    `toml:"max_batch_bytes"`
	TTLDuration           string `toml:"ttl-duration"`
	TTLNumBlocks          int    `toml:"ttl-num-blocks"`
}

type StateSyncConfig struct {
	Enable              bool   `toml:"enable"`
	RPCServers          string `toml:"rpc_servers"`
	TrustHeight         int    `toml:"trust_height"`
	TrustHash           string `toml:"trust_hash"`
	TrustPeriod         string `toml:"trust_period"`
	DiscoveryTime       string `toml:"discovery_time"`
	TempDir             string `toml:"temp_dir"`
	ChunkRequestTimeout string `toml:"chunk_request_timeout"`
	ChunkFetchers       string `toml:"chunk_fetchers"`
}

type BlockSyncConfig struct {
	Version string `toml:"version"`
}

type ConsensusConfig struct {
	WALFile                     string `toml:"wal_file"`
	TimeoutPropose              string `toml:"timeout_propose"`
	TimeoutProposeDelta         string `toml:"timeout_propose_delta"`
	TimeoutPrevote              string `toml:"timeout_prevote"`
	TimeoutPrevoteDelta         string `toml:"timeout_prevote_delta"`
	TimeoutPrecommit            string `toml:"timeout_precommit"`
	TimeoutPrecommitDelta       string `toml:"timeout_precommit_delta"`
	TimeoutCommit               string `toml:"timeout_commit"`
	DoubleSignCheckHeight       int    `toml:"double_sign_check_height"`
	SkipTimeoutCommit           bool   `toml:"skip_timeout_commit"`
	CreateEmptyBlocks           bool   `toml:"create_empty_blocks"`
	CreateEmptyBlocksInterval   string `toml:"create_empty_blocks_interval"`
	PeerGossipSleepDuration     string `toml:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration string `toml:"peer_query_maj23_sleep_duration"`
}

type StorageConfig struct {
	DiscardABCIResponses bool `toml:"discard_abci_responses"`
}

type TxIndexConfig struct {
	Indexer  string `toml:"indexer"`
	PSQLConn string `toml:"psql-conn"`
}

type InstrumentationConfig struct {
	Prometheus           bool   `toml:"prometheus"`
	PrometheusListenAddr string `toml:"prometheus_listen_addr"`
	MaxOpenConnections   int    `toml:"max_open_connections"`
	Namespace            string `toml:"namespace"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ProxyApp:               "tcp://127.0.0.1:26658",
		Moniker:                "KIRA VALIDATOR NODE",
		BlockSync:              true, // default value
		DBBackend:              "goleveldb",
		DBDir:                  "data",
		LogLevel:               "info",
		LogFormat:              "plain",
		GenesisFile:            "config/genesis.json",
		PrivValidatorKeyFile:   "config/priv_validator_key.json",
		PrivValidatorStateFile: "data/priv_validator_state.json",
		PrivValidatorLaddr:     "",
		NodeKeyFile:            "config/node_key.json",
		ABCI:                   "socket",
		FilterPeers:            false,
		//[blocksync]
		BlockSyncSection: BlockSyncConfig{
			Version: "v0",
		},

		// [rpc]
		RPC: RPCConfig{
			Laddr:                                "tcp://0.0.0.0:26657", // 0.0.0.0 or localhost?
			CORSAllowedOrigins:                   []string{"*"},
			CORSAllowedMethods:                   []string{"HEAD", "GET", "POST"},
			CORSAllowedHeaders:                   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"},
			GRPCLaddr:                            "",
			GRPCMaxOpenConnections:               900,
			Unsafe:                               false,
			MaxOpenConnections:                   900,
			MaxSubscriptionClients:               100,
			MaxSubscriptionsPerClient:            5,
			ExperimentalSubscriptionBufferSize:   200,
			ExperimentalWebsocketWriteBufferSize: 200,
			ExperimentalCloseOnSlowClient:        false,
			TimeoutBroadcastTxCommit:             "10s",
			MaxBodyBytes:                         1000000,
			MaxHeaderBytes:                       1048576,
			TLSCertFile:                          "",
			TLSKeyFile:                           "",
			PprofLaddr:                           "localhost:6060",
		},
		// [p2p]
		P2P: P2PConfig{
			Laddr:                        "tcp://0.0.0.0:26656",
			ExternalAddress:              "", // tcp://IP:PORT
			Seeds:                        "", // tcp://NODE_ID@IP:PORT
			PersistentPeers:              "",
			UPNP:                         false,
			AddrBookFile:                 "config/addrbook.json",
			AddrBookStrict:               false, //default is true
			MaxNumInboundPeers:           128,
			MaxNumOutboundPeers:          32,
			UnconditionalPeerIDs:         "",
			PersistentPeersMaxDialPeriod: "0s",
			FlushThrottleTimeout:         "100ms",
			MaxPacketMsgPayloadSize:      131072,
			SendRate:                     65536000,
			RecvRate:                     65536000,
			PEX:                          true,
			SeedMode:                     false,
			PrivatePeerIDs:               "",
			AllowDuplicateIP:             true,
			HandshakeTimeout:             "60s",
			DialTimeout:                  "30s",
		},

		// [mempool]
		Mempool: MempoolConfig{
			Version:               "v0",
			Recheck:               true,
			Broadcast:             true,
			WALDir:                "",
			Size:                  5000,
			MaxTxsBytes:           131072000,
			CacheSize:             10000,
			KeepInvalidTxsInCache: false,
			MaxTxBytes:            131072,
			MaxBatchBytes:         0,
			TTLDuration:           "0s",
			TTLNumBlocks:          0,
		},
		// [statesync]
		StateSync: StateSyncConfig{
			Enable:              false,
			RPCServers:          "",
			TrustHeight:         0,
			TrustHash:           "",
			TrustPeriod:         "168h0m0s",
			DiscoveryTime:       "15s",
			TempDir:             "/tmp",
			ChunkRequestTimeout: "10s",
			ChunkFetchers:       "4",
		},
		// [consensus]
		Consensus: ConsensusConfig{
			WALFile:                     "data/cs.wal/wal",
			TimeoutPropose:              "3s",
			TimeoutProposeDelta:         "500ms",
			TimeoutPrevote:              "1s",
			TimeoutPrevoteDelta:         "500ms",
			TimeoutPrecommit:            "1s",
			TimeoutPrecommitDelta:       "500ms",
			TimeoutCommit:               "3000ms",
			DoubleSignCheckHeight:       0,
			SkipTimeoutCommit:           false,
			CreateEmptyBlocks:           true,
			CreateEmptyBlocksInterval:   "20s",
			PeerGossipSleepDuration:     "100ms",
			PeerQueryMaj23SleepDuration: "2s",
		},
		// [storage]
		Storage: StorageConfig{
			DiscardABCIResponses: false,
		},
		// [tx_index]
		TxIndex: TxIndexConfig{
			Indexer:  "kv",
			PSQLConn: "",
		},
		// [instrumentation]
		Instrumentation: InstrumentationConfig{
			Prometheus:           true,
			PrometheusListenAddr: ":26660",
			MaxOpenConnections:   3,
			Namespace:            "tendermint",
		},
	}
}
