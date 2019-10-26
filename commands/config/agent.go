// Package config contain command to retrieve the server configuration
package config

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type agent struct {
	server string `short:"c" long:"ctm" description:"server" required:"true"`
	agent  string `short:"H" long:"host" description:"agent" required:"true"`
	all    bool   `short:"a" long:"all" description:"show all parameters, not only the default ones"`
	params types.ConfigAgentParamsReply
}

func (cmd *agent) Execute([]string) (err error) {
	if err = client.Call("GET", "/config/server/"+cmd.server+"/agent/"+cmd.agent+"/params", nil, map[string]string{}, &cmd.params); err != nil {
		return
	}
	return
}

func (cmd *agent) Data() interface{} {
	return cmd.params
}

func (cmd *agent) PrettyPrint() error {
	if cmd.all {
		for _, param := range cmd.params {
			defaultValue := ""
			if param.Value != param.DefaultValue {
				defaultValue = fmt.Sprintf("\t(default=%s)", param.DefaultValue)
			}
			fmt.Printf("%s=%s%s\n", param.Name, param.Value, defaultValue)
		}
	} else {
		for _, param := range cmd.params {
			if param.Value != param.DefaultValue {
				fmt.Printf("%s=%s\n", param.Name, param.Value)
			}
		}
	}
	return nil
}

func init() {
	commands.AddCommand("config.agent", "config.agent", "config.agent", &agent{})
}
