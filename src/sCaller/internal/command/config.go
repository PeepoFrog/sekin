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
	Overwrite string `json:"overwrite"`
	// Global flags:
	Home   string `json:"home"`
	LogFmt string `json:"log_format,omitempty"`
	LogLvl string `json:"log_level,omitempty"`
	Trace  string `json:"trace,omitempty"`
}
