// Package config contain command to retrieve the server configuration
package config

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type server struct {
	Server     string `short:"c" long:"ctm" description:"server" required:"true"`
	agents     types.ConfigAgentsReply
	hostgroups types.ConfigHostGroupsReply
}

func (cmd *server) Execute([]string) (err error) {
	if err = client.Call("GET", "/config/server/"+cmd.Server+"/agents", nil, map[string]string{}, &cmd.agents); err != nil {
		return err
	}
	_ = client.Call("GET", "/config/server/"+cmd.Server+"/hostgroups", nil, map[string]string{}, &cmd.hostgroups)
	return
}

func (cmd *server) Data() interface{} {
	return cmd
}

func (cmd *server) PrettyPrint() error {
	fmt.Printf("%50.50s %20.20s\n%s\n", "Node ID", "Status", strings.Repeat("-", 71))
	for _, agent := range cmd.agents.Agents {
		fmt.Printf("%50.50s %20.20s \n", agent.NodeID, agent.Status)
	}
	if cmd.hostgroups != nil && len(cmd.hostgroups) > 0 {
		fmt.Printf("\nHostgroups:\n%s\n", strings.Repeat("-", 71))
		for _, host := range cmd.hostgroups {
			fmt.Println(host)
		}
	}
	return nil
}

func init() {
	commands.AddCommand("config.server", "list the agents and hostgroups for a server", "List all the agent hostnames and all the hostgroups for a server", &server{})
}
