package job

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type JobsStatusCommand struct {
	application string
	limit       int
	status      string
	reply       types.JobsStatusReply
}

func (cmd *JobsStatusCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.application, "application", "", "Jobs for this application")
	flags.IntVar(&cmd.limit, "limit", 1000, "Limit to how many jobs")
	flags.StringVar(&cmd.status, "status", "", "Only this status")
}
func (cmd *JobsStatusCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	i = nil

	// add authorization header to the req
	args := make(map[string]string)
	if cmd.application != "" {
		args["application"] = cmd.application
	}
	if cmd.status != "" {
		args["status"] = cmd.status
	}
	args["limit"] = strconv.Itoa(cmd.limit)

	err = client.Call("GET", "/run/jobs/status", nil, args, &cmd.reply)

	i = cmd.reply

	return
}

func (cmd *JobsStatusCommand) PrettyPrint(f *flag.FlagSet, data interface{}) error {
	fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.16s %16.16s %5.5s %12.12s %12.12s\n",
		"Folder/Name", "Held", "JobId", "Order", "Status", "Host", "Del?", "Start time", "End time")
	fmt.Printf("-----------------------------------------------------------------------------------------------------------------\n")
	for _, job := range cmd.reply.Statuses {
		fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.146s %16.16s %5.5s %12.12s %12.12s\n",
			job.Folder+"/"+job.Name,
			strconv.FormatBool(job.Held),
			job.JobId, job.OrderDate, job.Status, job.Host, strconv.FormatBool(job.Deleted), job.StartTime, job.EndTime)
	}
	return nil
}
func init() {
	commands.Register("lj", &JobsStatusCommand{})
}
