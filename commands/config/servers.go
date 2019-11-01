// Package config contain command to retrieve the server configuration
package config

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type servers struct {
	servers types.ConfigServersReply
}

func (cmd *servers) Execute([]string) (err error) {
	err = client.Call("GET", "/config/servers", nil, map[string]string{}, &cmd.servers)
	return
}

func (cmd *servers) Data() interface{} {
	return cmd.servers
}

func (cmd *servers) PrettyPrint() error {
	fmt.Printf("%30.30s %10.10s %30.30s %10.10s\n%s\n", "Host", "Message", "Name", "State", strings.Repeat("-", 83))
	for _, server := range cmd.servers {
		fmt.Printf("%30.30s %10.10s %30.30s %10.10s\n", server.Host, server.Message, server.Name, server.State)
	}
	return nil
}

func init() {
	commands.AddCommand("config.servers", "list the servers", "List all the ctm servers", &servers{})
}
