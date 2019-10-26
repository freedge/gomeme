// Package config contain command to retrieve the server configuration
package config

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type agents struct {
	server string `short:"c" long:"ctm" description:"server" required:"true"`
	agents types.ConfigAgentsReply
}

func (cmd *agents) Execute([]string) (err error) {
	err = client.Call("GET", "/config/server/"+cmd.server+"/agents", nil, map[string]string{}, &cmd.agents)
	return
}

func (cmd *agents) Data() interface{} {
	return cmd.agents
}

func (cmd *agents) PrettyPrint() error {
	fmt.Printf("%20.20s %20.20s\n%s\n", "Node ID", "Status", strings.Repeat("-", 41))
	for _, agent := range cmd.agents.Agents {
		fmt.Printf("%20.20s %20.20s \n", agent.NodeID, agent.Status)
	}
	return nil
}

func init() {
	commands.AddCommand("config.agents", "config.agents", "config.agents", &agents{})
}
