// Package qr accesses Quantitative Resources
package qr

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type QRSetCommand struct {
	Name  string
	Ctm   string
	Max   int
	reply types.SetResourceReply
}

func (cmd *QRSetCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.Name, "name", "", "resource name")
	flags.StringVar(&cmd.Ctm, "ctm", "", "ctm")
	flags.IntVar(&cmd.Max, "max", -1, "max")
}
func (cmd *QRSetCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	if cmd.Name == "" || cmd.Ctm == "" || cmd.Max < 0 {
		err = fmt.Errorf("some argument is missing")
		return
	}
	err = client.Call("POST", "/run/resource/"+cmd.Ctm+"/"+cmd.Name, types.SetResourceQuery{Max: strconv.Itoa(cmd.Max)}, map[string]string{}, &cmd.reply)
	i = &cmd.reply

	return
}

func (cmd *QRSetCommand) PrettyPrint(f *flag.FlagSet, i interface{}) error {
	fmt.Println(cmd.reply.Message)
	return nil
}

func init() {
	commands.Register("qr.set", &QRSetCommand{})
}
