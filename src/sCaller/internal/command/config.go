package command

const ExecPath = "/sekaid"

type CommandHandler func(interface{}) (string, error)

var CommandMapping = map[string]struct {
	ArgsStruct interface{}
	Handler    CommandHandler
}{
	"init":    {ArgsStruct: SekaiInit{}, Handler: SekaiInitCmd},
	"version": {ArgsStruct: SekaiVersion{}, Handler: SekaiVersionCmd},
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
}
