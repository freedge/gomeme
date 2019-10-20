// Package qr accesses Quantitative Resources
package qr

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type QRCommand struct {
	Name string
	Ctm  string
	qrs  []types.QR
}

func (cmd *QRCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.Name, "name", "", "resource name")
	flags.StringVar(&cmd.Ctm, "ctm", "", "ctm")
}
func (cmd *QRCommand) Run() (i interface{}, err error) {
	i = nil

	args := make(map[string]string)

	if cmd.Name != "" {
		args["name"] = cmd.Name
	}
	if cmd.Ctm != "" {
		args["ctm"] = cmd.Ctm
	}
	err = client.Call("GET", "/run/resources", nil, args, &cmd.qrs)
	i = cmd.qrs

	return
}

func (cmd *QRCommand) PrettyPrint(i interface{}) error {
	fmt.Println("QR                       Ctm        Available      Max")
	fmt.Println("======================================================")
	for _, qr := range cmd.qrs {
		fmt.Printf("%-24.24s %-8.8s %11.11s %8d\n", qr.Name, qr.Ctm, qr.Available, qr.Max)
	}
	return nil
}

func init() {
	commands.Register("qr", &QRCommand{})
}
