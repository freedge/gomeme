package job

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type JobsStatusCommand struct {
	application string
	limit       int
	status      string
	reply       types.JobsStatusReply
	jobname     string
	jobid       string
	folder      string
	csv         bool
}

func (cmd *JobsStatusCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.application, "application", "", "Jobs for this application")
	flags.IntVar(&cmd.limit, "limit", 1000, "Limit to how many jobs")
	flags.StringVar(&cmd.status, "status", "", "Only this status")
	flags.StringVar(&cmd.jobname, "jobname", "", "Job name")
	flags.StringVar(&cmd.jobid, "jobid", "", "Jobid")
	flags.StringVar(&cmd.folder, "folder", "", "Folder")
	flags.BoolVar(&cmd.csv, "csv", false, "csv output")
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
	if cmd.jobname != "" {
		args["jobname"] = cmd.jobname
	}
	if cmd.jobid != "" {
		args["jobid"] = cmd.jobid
	}
	if cmd.folder != "" {
		args["folder"] = cmd.folder
	}
	args["limit"] = strconv.Itoa(cmd.limit)

	err = client.Call("GET", "/run/jobs/status", nil, args, &cmd.reply)

	i = cmd.reply

	return
}

// GetDuration returns the number of sconds a job has executed or executed so far
// or -1 if it has not ran yet
func GetDuration(job types.Status) (duration float64) {
	duration = -1
	if job.StartTime != "" {
		from, _ := types.ParseTime(job.StartTime)
		to := time.Now()
		if job.EndTime != "" {
			to, _ = types.ParseTime(job.EndTime)
		}
		duration = to.Sub(from).Seconds()
	}
	return
}

func GetDurationAsString(job types.Status) (duration string) {
	d := GetDuration(job)
	switch {
	case d < 0:
		duration = ""
		return
	case d < 120:
		duration = fmt.Sprintf("%ds", int(d))
	case d < 7200:
		duration = fmt.Sprintf("%dm", int(d/60))
	default:
		duration = fmt.Sprintf("%dh", int(d/3600))
	}
	if job.EndTime == "" {
		duration += " so far"
	}
	return
}

func (cmd *JobsStatusCommand) PrintCsv() error {
	fmt.Printf("folder,name,status,duration,starttime,endtime")
	for _, job := range cmd.reply.Statuses {
		fmt.Printf("%s,%s,%s,%f,%s,%s\n", job.Folder, job.Name, job.Status, GetDuration(job), job.StartTime, job.EndTime)
	}
	return nil
}

func (cmd *JobsStatusCommand) PrettyPrint(f *flag.FlagSet, data interface{}) error {
	if cmd.csv {
		return cmd.PrintCsv()
	}
	fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.16s %16.16s %5.5s %12.12s %12.12s %20.20s %8.8s\n",
		"Folder/Name", "Held", "JobId", "Order", "Status", "Host", "Del?", "Start time", "End time", "Description", "Duration")
	fmt.Printf("----------------------------------------------------------------------------------------------------------------------------------------------------------------------------\n")
	for _, job := range cmd.reply.Statuses {
		fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.146s %16.16s %5.5s %12.12s %12.12s %20.20s %s\n",
			job.Folder+"/"+job.Name,
			strconv.FormatBool(job.Held),
			job.JobId, job.OrderDate, job.Status, job.Host, strconv.FormatBool(job.Deleted), job.StartTime, job.EndTime, job.Description, GetDurationAsString(job))
	}
	return nil
}
func init() {
	commands.Register("lj", &JobsStatusCommand{})
}
