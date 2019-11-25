package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobActionCommand struct {
	Jobid  string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	Action string `short:"a" long:"action" description:"action to run" choice:"hold" choice:"free" choice:"confirm" choice:"delete" choice:"undelete" choice:"setToOk" choice:"runNow" choice:"kill" required:"true"`
	result types.JobActionReply
}

func (cmd *jobActionCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobActionCommand) Execute([]string) (err error) {
	err = client.Call("POST", "/run/job/"+cmd.Jobid+"/"+cmd.Action, nil, map[string]string{}, &cmd.result)
	return
}

func (cmd *jobActionCommand) PrettyPrint() error {
	fmt.Println(cmd.result.Message)
	return nil
}

func init() {
	commands.AddCommand("job.action", "action on a job", "Perform an action (confirm, hold, free, etc.) on a job", &jobActionCommand{})
}
