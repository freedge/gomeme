// Package config contain command to retrieve the server configuration
package config

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type agent struct {
	server string
	agent  string
	all    bool
	params types.ConfigAgentParamsReply
}

func (cmd *agent) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.server, "server", "", "Server to target")
	flags.StringVar(&cmd.agent, "agent", "", "Agent to target")
	flags.BoolVar(&cmd.all, "all", false, "Show all parameters, not only the non default ones")

}
func (cmd *agent) Run() (i interface{}, err error) {

	if cmd.server == "" || cmd.agent == "" {
		err = fmt.Errorf("server or agent not specified")
		return
	}
	if err = client.Call("GET", "/config/server/"+cmd.server+"/agent/"+cmd.agent+"/params", nil, map[string]string{}, &cmd.params); err != nil {
		return
	}

	i = cmd.params
	return
}

func (cmd *agent) PrettyPrint(data interface{}) error {
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
	commands.Register("config.agent", &agent{})
}
