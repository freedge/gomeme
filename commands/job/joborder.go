package job

import (
	"fmt"
	"time"

	"github.com/freedge/gomeme/client"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type orderJobCommand struct {
	DontHold bool   `short:"D" long:"donthold" description:"do not hold the job after submission"`
	Ctm      string `short:"c" long:"ctm" description:"server" required:"true"`
	Folder   string `short:"f" long:"folder" description:"folder" required:"true"`
	Jobs     string `short:"n" long:"name" description:"job name" required:"true"`
	Retries  int    `short:"r" long:"retries" description:"try to get the created job id a few times" default:"2"`
	reply    types.OrderJobReply
	status   types.JobsStatusReply
}

func (cmd *orderJobCommand) Data() interface{} {
	return cmd.status
}
func (cmd *orderJobCommand) Execute([]string) (err error) {
	if err := commands.RequiresAnnotation(); err != nil {
		return err
	}
	query := types.OrderQuery{Jobs: cmd.Jobs, Ctm: cmd.Ctm, Folder: cmd.Folder, Hold: !cmd.DontHold}
	err = client.Call("POST", "/run/order", query, map[string]string{}, &cmd.reply)
	if err != nil {
		return
	}

	_ = client.Call("GET", "/run/status/"+cmd.reply.RunID, nil, map[string]string{}, &cmd.status)
	for cmd.status.Statuses == nil && cmd.Retries > 0 {
		cmd.Retries--
		time.Sleep(1 * time.Second)
		_ = client.Call("GET", "/run/status/"+cmd.reply.RunID, nil, map[string]string{}, &cmd.status)
	}

	return
}

func (cmd *orderJobCommand) PrettyPrint() error {
	fmt.Println("RunId: ", cmd.reply.RunID)
	if cmd.status.Statuses != nil {
		for _, status := range cmd.status.Statuses {
			fmt.Println("JobId: ", status.JobId)
		}
	} else {
		fmt.Println(cmd.reply.MonitorPageURI, cmd.reply.StatusURI)
	}
	return nil
}

func init() {
	commands.AddCommand("job.order", "order a job", "Order the specified job", &orderJobCommand{})
}
