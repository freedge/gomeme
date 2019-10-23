package login

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type logout struct {
	out types.LogoutReply
}

func (cmd *logout) Prepare(flags *flag.FlagSet) {}
func (cmd *logout) Run() (i interface{}, err error) {

	if err = client.Call("POST", "/session/logout", nil, map[string]string{}, &cmd.out); err != nil {
		return
	}

	i = cmd.out
	return
}

func (cmd *logout) PrettyPrint(data interface{}) error {
	fmt.Println(cmd.out.Message)
	return nil
}

func init() {
	commands.Register("logout", &logout{})
}
