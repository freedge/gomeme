package login

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type logout struct {
	out types.LogoutReply
}

func (cmd *logout) Data() interface{} { return cmd.out }
func (cmd *logout) Execute([]string) (err error) {
	err = client.Call("POST", "/session/logout", nil, map[string]string{}, &cmd.out)
	delete(commands.Tokens.Endpoint, commands.Opts.Endpoint)
	commands.WriteTokensFile()
	return
}

func (cmd *logout) PrettyPrint() error {
	fmt.Println(cmd.out.Message)
	return nil
}

func init() {
	commands.AddCommand("logout", "logout", "Logout the current token", &logout{})
}
