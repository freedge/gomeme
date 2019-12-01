package job

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type jobStatusCommand struct {
	Jobid         string `short:"j" long:"jobid" required:"true" description:"Job ID"`
	Result        types.Status
	WaitingInfo   types.WaitinfoReply
	JobDefinition types.JobGetReply
}

func (cmd *jobStatusCommand) Data() interface{} {
	return cmd
}

const defaultJobGetAnnotation = defaultJobLogAnnotation

func (cmd *jobStatusCommand) Execute([]string) (err error) {
	if commands.Opts.Subject == "" {
		// gomeme provide a default annotation for job get
		commands.Opts.Subject = defaultJobGetAnnotation
	}

	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/status", nil, map[string]string{}, &cmd.Result)

	if err == nil {
		_ = client.Call("GET", "/run/job/"+cmd.Jobid+"/waitingInfo", nil, map[string]string{}, &cmd.WaitingInfo)
		// not sure how to parse this at this moment

		err = client.Call("GET", "/run/job/"+cmd.Jobid+"/get", nil, map[string]string{}, &cmd.JobDefinition)

	}

	return
}

func (cmd *jobStatusCommand) PrettyPrint() error {
	fmt.Printf("%s/%s %s held=%v\n", cmd.Result.Folder, cmd.Result.Name, cmd.Result.Status, cmd.Result.Held)
	if cmd.WaitingInfo != nil {
		fmt.Println(cmd.WaitingInfo)
	}
	if cmd.JobDefinition != nil {
		for _, def := range cmd.JobDefinition {
			fmt.Printf("(runas %s) %s/%s %s\n", def.RunAs, def.FilePath, def.FileName, strings.Join(def.Arguments, " "))
		}
	}

	return nil
}

func init() {
	commands.AddCommand("job.get", "status of a jobid", "Retrieves the status of one job, hoping one day to have estimation of its completion", &jobStatusCommand{})
}
