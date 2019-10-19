package job

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type JobActionCommand struct {
	Jobid  string
	action string
	result types.JobActionReply
}

func (cmd *JobActionCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.Jobid, "jobid", "", "JobID")
	flags.StringVar(&cmd.action, "action", "", "action to run: hold, free, confirm, delete, undelete, rerun, setToOk")
}

func (cmd *JobActionCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	i = nil

	if cmd.Jobid == "" || cmd.action == "" {
		err = fmt.Errorf("job id or action missing")
		return
	}

	err = client.Call("POST", "/run/job/"+cmd.Jobid+"/"+cmd.action, nil, map[string]string{}, &cmd.result)
	i = cmd.result

	return
}

func (cmd *JobActionCommand) PrettyPrint(f *flag.FlagSet, i interface{}) error {
	fmt.Println(cmd.result.Message)
	return nil
}

func init() {
	commands.Register("job.action", &JobActionCommand{})
}
