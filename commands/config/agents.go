// Package config contain command to retrieve the server configuration
package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type agents struct {
	server string
	agents types.ConfigAgentsReply
}

func (cmd *agents) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.server, "server", "", "Server to target")
}
func (cmd *agents) Run() (i interface{}, err error) {

	if cmd.server == "" {
		err = fmt.Errorf("server not specified")
		return
	}
	if err = client.Call("GET", "/config/server/"+cmd.server+"/agents", nil, map[string]string{}, &cmd.agents); err != nil {
		return
	}

	i = cmd.agents
	return
}

func (cmd *agents) PrettyPrint(data interface{}) error {
	fmt.Printf("%20.20s %20.20s\n%s\n", "Node ID", "Status", strings.Repeat("-", 41))
	for _, agent := range cmd.agents.Agents {
		fmt.Printf("%20.20s %20.20s \n", agent.NodeID, agent.Status)
	}
	return nil
}

func init() {
	commands.Register("config.agents", &agents{})
}
