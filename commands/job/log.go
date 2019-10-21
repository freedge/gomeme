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
}

func (cmd *jobLogCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.jobid, "jobid", "", "JobID")
	flags.BoolVar(&cmd.output, "output", false, "display output instead of logs")
}
func (cmd *jobLogCommand) Run() (i interface{}, err error) {
	i = nil

	service := "log"
	if cmd.output {
		service = "output"
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
