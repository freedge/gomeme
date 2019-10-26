// Package qr accesses Quantitative Resources
package qr

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type listQRCommand struct {
	Name string `short:"n" long:"name" description:"resource name"`
	Ctm  string `short:"c" long:"ctm" description:"server"`
	qrs  []types.QR
}

func (cmd *listQRCommand) Data() interface{} {
	return cmd.qrs
}

func (cmd *listQRCommand) Execute([]string) (err error) {

	args := make(map[string]string)
	commands.AddIfNotEmpty(args, "name", cmd.Name)
	commands.AddIfNotEmpty(args, "ctm", cmd.Ctm)

	err = client.Call("GET", "/run/resources", nil, args, &cmd.qrs)
	return
}

func (cmd *listQRCommand) PrettyPrint() error {
	fmt.Println("QR                       Ctm        Available      Max")
	fmt.Println("======================================================")
	for _, qr := range cmd.qrs {
		fmt.Printf("%-24.24s %-8.8s %11.11s %8d\n", qr.Name, qr.Ctm, qr.Available, qr.Max)
	}
	return nil
}

func init() {
	commands.AddCommand("qr", "list qr", "List quantitative resources, optionally filtering them", &listQRCommand{})
}
