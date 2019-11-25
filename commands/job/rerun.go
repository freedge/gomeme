package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobRerunCommand struct {
	Jobid  string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	result types.Status
}

func (cmd *jobRerunCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobRerunCommand) Execute([]string) (err error) {
	err = client.Call("POST", "/run/job/"+cmd.Jobid+"/rerun", nil, map[string]string{}, &cmd.result)
	return
}

func (cmd *jobRerunCommand) PrettyPrint() error {
	fmt.Printf("%s/%s %s (%d runs) in status %s\n", cmd.result.Folder, cmd.result.Name, cmd.result.JobId, cmd.result.NumberOfRuns, cmd.result.Status)
	return nil
}

func init() {
	commands.AddCommand("job.rerun", "rerun a job", "Rerun a job", &jobRerunCommand{})
}
