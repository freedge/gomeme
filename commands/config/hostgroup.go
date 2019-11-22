// Package config contain command to retrieve the server configuration
package config

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type hostgroup struct {
	Server    string `short:"c" long:"ctm" description:"server" required:"true"`
	HostGroup string `short:"g" long:"hostgroup" description:"host group" required:"true"`

	reply types.HostGroupAgentsReply
}

func (cmd *hostgroup) Execute([]string) (err error) {
	err = client.Call("GET", "/config/server/"+cmd.Server+"/hostgroup/"+cmd.HostGroup+"/agents", nil, map[string]string{}, &cmd.reply)
	return
}

func (cmd *hostgroup) Data() interface{} {
	return cmd.reply
}

func (cmd *hostgroup) PrettyPrint() error {
	for _, agent := range cmd.reply {
		fmt.Printf("%s \n", agent.Host)
	}
	return nil
}

func init() {
	commands.AddCommand("config.hostgroup", "list the agents of a hostgroup", "List all the agent hostnames for a hostgroup", &hostgroup{})
}
