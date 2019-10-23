package job

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type orderJobCommand struct {
	Hold   bool
	Ctm    string
	Folder string
	Jobs   string
	reply  types.OrderJobReply
	status types.JobsStatusReply
}

func (cmd *orderJobCommand) Prepare(flags *flag.FlagSet) {
	flags.BoolVar(&cmd.Hold, "hold", true, "Hold the job after submission")
	flags.StringVar(&cmd.Ctm, "ctm", "", "ctm")
	flags.StringVar(&cmd.Folder, "folder", "", "Folder")
	flags.StringVar(&cmd.Jobs, "jobs", "", "jobs")
}
func (cmd *orderJobCommand) Run() (i interface{}, err error) {
	i = nil
	if commands.TheToken == "" {
		err = fmt.Errorf("no token found. Please login first")
		return
	}
	if !cmd.Hold {
		err = fmt.Errorf("Currently only support holding jobs")
		return
	}

	if cmd.Jobs == "" || cmd.Ctm == "" || cmd.Folder == "" {
		err = fmt.Errorf("parameters missing")
		return
	}

	query := types.OrderQuery{Jobs: cmd.Jobs, Ctm: cmd.Ctm, Folder: cmd.Folder, Hold: cmd.Hold}
	err = client.Call("POST", "/run/order", query, map[string]string{}, &cmd.reply)
	if err != nil {
		return
	}

	err = client.Call("GET", "/run/status/"+cmd.reply.RunID, nil, map[string]string{}, &cmd.status)
	i = &cmd.status

	return
}

func (cmd *orderJobCommand) PrettyPrint(data interface{}) error {
	fmt.Println("RunId: ", cmd.reply.RunID)
	for _, status := range cmd.status.Statuses {
		fmt.Println("JobId: ", status.JobId)
	}
	return nil
}

func init() {
	commands.Register("job.order", &orderJobCommand{})
}
