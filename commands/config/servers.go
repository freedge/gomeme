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

type servers struct {
	servers types.ConfigServersReply
}

func (cmd *servers) Prepare(flags *flag.FlagSet) {}
func (cmd *servers) Run() (i interface{}, err error) {

	if err = client.Call("GET", "/config/servers", nil, map[string]string{}, &cmd.servers); err != nil {
		return
	}

	i = cmd.servers
	return
}

func (cmd *servers) PrettyPrint(data interface{}) error {
	fmt.Printf("%10.10s %10.10s %10.10s %10.10s\n%s\n", "Host", "Message", "Name", "State", strings.Repeat("-", 43))
	for _, server := range cmd.servers {
		fmt.Printf("%10.10s %10.10s %10.10s %10.10s\n", server.Host, server.Message, server.Name, server.State)
	}
	return nil
}

func init() {
	commands.Register("config.servers", &servers{})
}
