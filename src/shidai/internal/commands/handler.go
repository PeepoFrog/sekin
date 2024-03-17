package commands

type CommandHandler interface {
	HandleCommand(args map[string]interface{}) error
}
