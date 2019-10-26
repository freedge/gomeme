// Package deploy call the deploy API
package deploy

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type get struct {
	Ctm    string `short:"c" long:"ctm" description:"server" required:"true"`
	Folder string `short:"f" long:"folder" description:"folder" required:"true"`
	Xml    bool   `short:"x" long:"xml" description:"xml format"`
	output string
}

func (cmd *get) Data() interface{} {
	return cmd.output
}

func (cmd *get) Execute([]string) (err error) {
	params := map[string]string{"ctm": cmd.Ctm, "folder": cmd.Folder}
	if cmd.Xml {
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
	commands.AddCommand("deploy.get", "get jobs definition", "Get all the jobs definition under a folder", &get{})
}
