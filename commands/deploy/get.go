// Package deploy call the deploy API
package deploy

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type get struct {
	ctm    string `short:"c" long:"ctm" description:"server" required:"true"`
	folder string `short:"f" long:"folder" description:"folder" required:"true"`
	xml    bool   `short:"x" long:"xml" description:"xml format"`
	output string
}

func (cmd *get) Data() interface{} {
	return cmd.output
}

func (cmd *get) Execute([]string) (err error) {
	params := map[string]string{"ctm": cmd.ctm, "folder": cmd.folder}
	if cmd.xml {
		params["format"] = "XML"
	}

	err = client.Call("GET", "/deploy/jobs", nil, params, &cmd.output)
	return
}

func (cmd *get) PrettyPrint() error {
	fmt.Println(cmd.output)
	return nil
}

func init() {
	commands.AddCommand("deploy.get", "deploy.get", "deploy.get", &get{})
}
