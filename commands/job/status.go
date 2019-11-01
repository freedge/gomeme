package job

import (
	"encoding/json"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobStatusCommand struct {
	Jobid       string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	Result      types.Status
	WaitingInfo string
}

func (cmd *jobStatusCommand) Data() interface{} {
	return cmd
}

func (cmd *jobStatusCommand) Execute([]string) (err error) {
	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/status", nil, map[string]string{}, &cmd.Result)

	if err != nil {
		_ = client.Call("GET", "/run/job/"+cmd.Jobid+"/waitingInfo", nil, map[string]string{}, &cmd.WaitingInfo)
		// not sure how to parse this at this moment
	}

	return
}

func (cmd *jobStatusCommand) PrettyPrint() error {
	// this command is useless at this moment, just display the json as parsed
	bytes, err := json.MarshalIndent(cmd, "", "  ")
	if err == nil {
		fmt.Printf("%s\n", string(bytes))
	} else {
		fmt.Printf("%#v %#v %#v\n", cmd.Result, cmd.WaitingInfo, err)
	}
	return nil
}

func init() {
	commands.AddCommand("job.get", "status of a jobid", "Retrieves the status of one job, hoping one day to have estimation of its completion", &jobStatusCommand{})
}
