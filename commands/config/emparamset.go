package config

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

// only there to bootstrap our local workbench.
// Maybe rename as bootstrap instead?
type emparamset struct {
	reply      string
	ParamName  string `long:"name" required:"true" choice:"UserAuditAnnotationOn" description:"the parameter to change"`
	ParamValue string `long:"value" required:"true"`
}

type emparam struct {
	Value string `json:"value"`
}

func (cmd *emparamset) Execute([]string) (err error) {
	err = client.Call("POST", "/config/em/param/"+cmd.ParamName, emparam{cmd.ParamValue}, map[string]string{}, &cmd.reply)
	return
}

func (cmd *emparamset) Data() interface{} {
	return cmd.reply
}

func (cmd *emparamset) PrettyPrint() error {
	fmt.Printf("%s\n", cmd.reply)
	return nil
}

func init() {
	commands.AddCommand("test.config.emparamset", "set an emparameter", "Only used to bootstrap our workbench", &emparamset{})
}
