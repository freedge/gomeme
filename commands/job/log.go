package job

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type jobLogCommand struct {
	jobid  string `short:"j" long:"jobid" description:"jobid" required:"true"`
	output bool   `short:"o" long:"output"  description:"display output instead of logs"`
	result string
	run    int `short:"r" long:"run" default:"-1" description:"for output, run number of the job. Defaults to last one"`
}

func (cmd *jobLogCommand) Data() interface{} {
	return cmd.result
}

func (cmd *jobLogCommand) Execute([]string) (err error) {
	service := "log"
	if cmd.output {
		service = "output"
	}
	if cmd.run > 0 {
		if cmd.output {
			service = fmt.Sprintf("%s/?runNo=%d", service, cmd.run)
		} else {
			err = fmt.Errorf("run can only be specified when getting a job output")
			return
		}
	}

	err = client.Call("GET", "/run/job/"+cmd.jobid+"/"+service, nil, map[string]string{}, &cmd.result)
	return
}

func (cmd *jobLogCommand) PrettyPrint() error {
	fmt.Println(cmd.result)
	return nil
}

func init() {
	commands.AddCommand("job.log", "job.log", "job.log", &jobLogCommand{})
}
