package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobActionCommand struct {
	jobid  string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	action string `short:"a" long:"action" description:"action to run" choice:"hold" choice:"free" choice:"confirm" choice:"delete" choice:"undelete" choice:"rerun" choice:"setToOk" choice:"runNow" choice:"kill" required:"true"`
	result types.JobActionReply
}

func (cmd *jobActionCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobActionCommand) Execute([]string) (err error) {
	err = client.Call("POST", "/run/job/"+cmd.jobid+"/"+cmd.action, nil, map[string]string{}, &cmd.result)
	return
}

func (cmd *jobActionCommand) PrettyPrint() error {
	fmt.Println(cmd.result.Message)
	return nil
}

func init() {
	commands.AddCommand("job.action", "job.action", "job.action", &jobActionCommand{})
}
