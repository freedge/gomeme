package job

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type jobLogCommand struct {
	jobid  string
	output bool
	result string
	run    int
}

func (cmd *jobLogCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.jobid, "jobid", "", "JobID")
	flags.BoolVar(&cmd.output, "output", false, "display output instead of logs")
	flags.IntVar(&cmd.run, "run", -1, "for output, run number of the job. Defaults to last one")
}
func (cmd *jobLogCommand) Run() (i interface{}, err error) {
	i = nil

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
	i = cmd.result

	return
}

func (cmd *jobLogCommand) PrettyPrint(i interface{}) error {
	fmt.Println(cmd.result)
	return nil
}

func init() {
	commands.Register("job.log", &jobLogCommand{})
}
