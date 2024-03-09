package command

const ExecPath = "/sekaid"

type CommandHandler func(interface{}) (string, error)

var CommandMapping = map[string]struct {
	ArgsStruct interface{}
	Handler    CommandHandler
}{
	"init":                {ArgsStruct: SekaiInit{}, Handler: SekaiInitCmd},
	"version":             {ArgsStruct: SekaiVersion{}, Handler: SekaiVersionCmd},
	"add-genesis-account": {ArgsStruct: SekaiAddGenesisAcc{}, Handler: SekaiAddGenesisAccCmd},
	"gentx-claim":         {ArgsStruct: SekaiGentxClaim{}, Handler: SekaiGentxClaimCmd},
}

type SekaiVersion struct {
}

type SekaiInit struct {
	ChainID   string `json:"chain-id"`
	Overwrite bool   `json:"overwrite"`
	Moniker   string `json:"moniker"`
	// Global flags:
	Home   string `json:"home"`
	LogFmt string `json:"log_format,omitempty"`
	LogLvl string `json:"log_level,omitempty"`
	Trace  string `json:"trace,omitempty"`
}

<<<<<<< HEAD
type SekaidKeysAdd struct {
	KeyName        string `json:"key-name"`
	KeyringBackend string `json:"keyring-backend"`
	Recover        bool   `json:"recover"`
	// Global flags:
	Home   string `json:"home"`
	LogFmt string `json:"log_format,omitempty"`
	LogLvl string `json:"log_level,omitempty"`
	Trace  string `json:"trace,omitempty"`
	Output string `json:"output,omitempty"`
=======
type SekaiAddGenesisAcc struct {
	Address string   `json:"address"` // Key can be used instead of address
	Coins   []string `json:"coins"`
	// Global flags:
	Home    string `json:"home"`
	Keyring string `json:"keyring-backend"`
	LogFmt  string `json:"log_format,omitempty"`
	LogLvl  string `json:"log_level,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

type SekaiGentxClaim struct {
	Address string `json:"address"`
	Keyring string `json:"keyring-backend"`
	Moniker string `json:"moniker"`
	PubKey  string `json:"pubkey"`
	Home    string `json:"home"`
	LogFmt  string `json:"log_format,omitempty"`
	LogLvl  string `json:"log_level,omitempty"`
	Trace   bool   `json:"trace,omitempty"`
>>>>>>> 6ff97d6 (feat(rest_api) Add add-genesis-account, gentx-claim)
}
