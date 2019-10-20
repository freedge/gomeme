package job

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type JobLogCommand struct {
	Jobid  string
	Output bool
	result string
}

func (cmd *JobLogCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.Jobid, "jobid", "", "JobID")
	flags.BoolVar(&cmd.Output, "output", false, "display output instead of logs")
}
func (cmd *JobLogCommand) Run() (i interface{}, err error) {
	i = nil

	service := "log"
	if cmd.Output {
		service = "output"
	}

	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/"+service, nil, map[string]string{}, &cmd.result)
	i = cmd.result

	return
}

func (cmd *JobLogCommand) PrettyPrint(i interface{}) error {
	fmt.Println(cmd.result)
	return nil
}

func init() {
	commands.Register("job.log", &JobLogCommand{})
}
