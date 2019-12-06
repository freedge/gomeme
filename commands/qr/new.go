package qr

import (
	"strconv"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

// only there to bootstrap our workbench
type newQRCommand struct {
	setQRCommand
}

func (cmd *newQRCommand) Execute([]string) (err error) {
	if err := commands.RequiresAnnotation(); err != nil {
		return err
	}
	err = client.Call("POST", "/run/resource/"+cmd.Ctm, types.AddResourceQuery{Max: strconv.Itoa(cmd.Max), Name: cmd.Name}, map[string]string{}, &cmd.reply)
	return
}

func init() {
	commands.AddCommand("test.qr.new", "Create a new qr", "Create a new QR, only used to bootstrap our workbench", &newQRCommand{})
}
