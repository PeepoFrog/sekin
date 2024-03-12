package command

type (
	Handler func(interface{}) (string, error)
	//example
	// SekaiInit struct {
	// 	ChainID   string `json:"chain-id"`
	// 	Moniker   string `json:"moniker"`
	// 	Home      string `json:"home"`
	// 	LogFmt    string `json:"log_format"`
	// 	LogLvl    string `json:"log_level"`
	// 	Trace     bool   `json:"trace"`
	// 	Overwrite bool   `json:"overwrite"`
	// }

	InitNew struct {
		Mnemonic string `json:"mnemonic"`
	}
)

var CommandMapping = map[string]struct {
	ArgsStruct func() interface{}
	Handler    Handler
}{
	// example "init":                {ArgsStruct: func() interface{} { return &SekaiInit{} }, Handler: SekaiInitCmd},
	"init-new": {ArgsStruct: func() interface{} { return &InitNew{} }, Handler: InitNewNetwork},
}
