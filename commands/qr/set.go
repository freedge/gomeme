package qr

import (
	"fmt"
	"strconv"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type setQRCommand struct {
	Name  string `short:"n" long:"name" required:"true" description:"resource name"`
	Ctm   string `short:"c" long:"ctm" required:"true" description:"ctm"`
	Max   int    `short:"m" long:"max" required:"true" description:"max number of QRs"`
	reply types.SetResourceReply
}

func (cmd *setQRCommand) Data() interface{} {
	return cmd.reply
}

func (cmd *setQRCommand) Execute([]string) (err error) {
	err = client.Call("POST", "/run/resource/"+cmd.Ctm+"/"+cmd.Name, types.SetResourceQuery{Max: strconv.Itoa(cmd.Max)}, map[string]string{}, &cmd.reply)
	return
}

func (cmd *setQRCommand) PrettyPrint() error {
	fmt.Println(cmd.reply.Message)
	return nil
}

func init() {
	commands.AddCommand("qr.set", "set a qr", "Set a quantitative resource under a given ctm server", &setQRCommand{})
}
