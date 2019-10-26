package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type orderJobCommand struct {
	DontHold bool   `short:"D" long:"donthold" description:"do not hold the job after submission"`
	Ctm      string `short:"c" long:"ctm" description:"server" required:"true"`
	Folder   string `short:"f" long:"folder" description:"folder" required:"true"`
	Jobs     string `short:"n" long:"name" description:"job name" required:"true"`
	reply    types.OrderJobReply
	status   types.JobsStatusReply
}

func (cmd *orderJobCommand) Data() interface{} {
	return cmd.status
}
func (cmd *orderJobCommand) Execute([]string) (err error) {
	query := types.OrderQuery{Jobs: cmd.Jobs, Ctm: cmd.Ctm, Folder: cmd.Folder, Hold: !cmd.DontHold}
	err = client.Call("POST", "/run/order", query, map[string]string{}, &cmd.reply)
	if err != nil {
		return
	}

	err = client.Call("GET", "/run/status/"+cmd.reply.RunID, nil, map[string]string{}, &cmd.status)
	return
}

func (cmd *orderJobCommand) PrettyPrint() error {
	fmt.Println("RunId: ", cmd.reply.RunID)
	for _, status := range cmd.status.Statuses {
		fmt.Println("JobId: ", status.JobId)
	}
	return nil
}

func init() {
	commands.AddCommand("job.order", "job.order", "job.order", &orderJobCommand{})
}
