// Package deploy call the deploy API
package deploy

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type get struct {
	ctm    string
	folder string
	xml    bool
	output string
}

func (cmd *get) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.ctm, "ctm", "", "ctm")
	flags.StringVar(&cmd.folder, "folder", "", "folder")
	flags.BoolVar(&cmd.xml, "xml", false, "xml format")
}

func (cmd *get) Run() (i interface{}, err error) {

	if cmd.folder == "" || cmd.ctm == "" {
		err = fmt.Errorf("folder or ctm not specified")
		return
	}
	params := map[string]string{"ctm": cmd.ctm, "folder": cmd.folder}
	if cmd.xml {
		params["format"] = "XML"
	}

	if err = client.Call("GET", "/deploy/jobs", nil, params, &cmd.output); err != nil {
		return
	}

	i = cmd.output
	return
}

func (cmd *get) PrettyPrint(data interface{}) error {
	fmt.Println(cmd.output)
	return nil
}

func init() {
	commands.Register("deploy.get", &get{})
}
