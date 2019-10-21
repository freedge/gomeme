package job

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobActionCommand struct {
	jobid  string
	action string
	result types.JobActionReply
}

func (cmd *jobActionCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.jobid, "jobid", "", "JobID")
	flags.StringVar(&cmd.action, "action", "", "action to run: hold, free, confirm, delete, undelete, rerun, setToOk")
}

func (cmd *jobActionCommand) Run() (i interface{}, err error) {
	i = nil

	if cmd.jobid == "" || cmd.action == "" {
		err = fmt.Errorf("job id or action missing")
		return
	}

	err = client.Call("POST", "/run/job/"+cmd.jobid+"/"+cmd.action, nil, map[string]string{}, &cmd.result)
	i = cmd.result

	return
}

func (cmd *jobActionCommand) PrettyPrint(i interface{}) error {
	fmt.Println(cmd.result.Message)
	return nil
}

func init() {
	commands.Register("job.action", &jobActionCommand{})
}
