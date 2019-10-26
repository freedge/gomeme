package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type jobLogCommand struct {
	Jobid  string `short:"j" long:"jobid" description:"jobid" required:"true"`
	Output bool   `short:"o" long:"output"  description:"display output instead of logs"`
	result string
	Run    int `short:"r" long:"run" default:"-1" description:"for output, run number of the job. Defaults to last one"`
}

func (cmd *jobLogCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobLogCommand) Execute([]string) (err error) {
	service := "log"
	if cmd.Output {
		service = "output"
	}
	if cmd.Run > 0 {
		if cmd.Output {
			service = fmt.Sprintf("%s/?runNo=%d", service, cmd.Run)
		} else {
			err = fmt.Errorf("run can only be specified when getting a job output")
			return
		}
	}

	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/"+service, nil, map[string]string{}, &cmd.result)
	return
}

func (cmd *jobLogCommand) PrettyPrint() error {
	fmt.Println(cmd.result)
	return nil
}

func init() {
	commands.AddCommand("job.log", "get logs for a job", "Retrieve output or logs of a job id", &jobLogCommand{})
}
