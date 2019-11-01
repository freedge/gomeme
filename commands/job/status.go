package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobStatusCommand struct {
	Jobid       string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	result      types.Status
	waitingInfo string
}

func (cmd *jobStatusCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobStatusCommand) Execute([]string) (err error) {
	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/status", nil, map[string]string{}, &cmd.result)

	if err != nil {
		_ = client.Call("GET", "/run/job/"+cmd.Jobid+"/waitingInfo", nil, map[string]string{}, &cmd.result)
		// not sure how to parse this at this moment
	}

	return
}

func (cmd *jobStatusCommand) PrettyPrint() error {
	fmt.Printf("%#v %#v\n", cmd.result, cmd.waitingInfo)
	return nil
}

func init() {
	commands.AddCommand("job.get", "status of a jobid", "Retrieves the status of one job, hoping one day to have estimation of its completion", &jobStatusCommand{})
}
