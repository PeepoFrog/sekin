package commands

type JoinCommandHandler struct{}

func (j *JoinCommandHandler) HandleCommand(args map[string]interface{}) error {
	return nil
}

func init() {
	RegisterCommand("join", &JoinCommandHandler{})
}
