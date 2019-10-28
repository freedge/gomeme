package config

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type ping struct {
	Server   string `short:"c" long:"ctm" description:"server" required:"true"`
	Agent    string `short:"H" long:"host" description:"agent" required:"true"`
	Discover bool   `short:"d" long:"discover" description:"discover"`
	Timeout  int    `short:"t" long:"timeout" default:"10" description:"timeout"`
	params   types.PingAgentReply
}

func (cmd *ping) Execute([]string) (err error) {
	query := types.PingAgentQuery{Timeout: cmd.Timeout, Discover: cmd.Discover}
	err = client.Call("POST", "/config/server/"+cmd.Server+"/agent/"+cmd.Agent+"/ping", query, map[string]string{}, &cmd.params)
	return
}

func (cmd *ping) Data() interface{} {
	return cmd.params.Message
}

func (cmd *ping) PrettyPrint() error {
	fmt.Println(cmd.params.Message)
	return nil
}

func init() {
	commands.AddCommand("config.ping", "ping an agent", "Ping an agent", &ping{})
}
