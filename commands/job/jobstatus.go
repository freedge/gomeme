package job

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

// jobsStatusCommand retrieve a list of jobs
type jobsStatusCommonCommand struct {
	Application string `short:"a" long:"application"`
	Limit       int    `short:"l"  long:"limit" default:"1000"`
	Status      string `short:"s"  long:"status" choice:"Executing" choice:"Ended Not OK" choice:"Ended OK" choice:"Wait Condition" choice:"Wait Resource" choice:"Wait User" choice:"Wait Host" description:"Only this status"`
	reply       types.JobsStatusReply
	Jobname     string `short:"n" long:"jobname" description:"job name"`
	Jobid       string `short:"j" long:"jobid" description:"job id"`
	Folder      string `short:"f" long:"folder" description:"folder"`
	Verbose     bool   `short:"v" long:"verbose" description:"output more stuff, that fits onto my sreen (158 characters wide)"`
	Host        string `short:"H" long:"host" description:"host"`
	Neighbours  bool   `long:"deps" description:"browse through neighours of this job. Only jobid can be used to filter jobs"`
	Ctm         string `short:"c" long:"ctm" description:"Control-m server to target"`
}

type jobsStatusCommand struct {
	jobsStatusCommonCommand
	Csv bool `short:"c" long:"csv" description:"csv output"`
}

func (cmd *jobsStatusCommand) Data() interface{} {
	return cmd.reply
}

const (
	// URL to target to get the jobs status
	jobsStatusPath = "/run/jobs/status"
)

// func addArg(args map[string]string, )

func (cmd *jobsStatusCommonCommand) GetJobs() (i interface{}, err error) {
	i = nil

	// add authorization header to the req
	args := make(map[string]string)
	commands.AddIfNotEmpty(args, "application", cmd.Application)
	commands.AddIfNotEmpty(args, "status", cmd.Status)
	commands.AddIfNotEmpty(args, "jobname", cmd.Jobname)
	commands.AddIfNotEmpty(args, "folder", cmd.Folder)
	commands.AddIfNotEmpty(args, "host", cmd.Host)
	commands.AddIfNotEmpty(args, "ctm", cmd.Ctm)

	if cmd.Jobid != "" {
		if cmd.Neighbours {
			if len(args) > 0 {
				err = fmt.Errorf("only jobid should be used to filter jobs")
				return
			}
			args["neighborhood"] = "1"
			args["direction"] = "radial"
			args["depth"] = "5"
		}
		args["jobid"] = cmd.Jobid
	} else {
		if cmd.Neighbours {
			err = fmt.Errorf("jobid missing")
			return
		}
	}
	args["limit"] = strconv.Itoa(cmd.Limit)

	err = client.Call("GET", jobsStatusPath, nil, args, &cmd.reply)

	i = cmd.reply

	return
}

func (cmd *jobsStatusCommand) Execute([]string) (err error) {
	_, err = cmd.GetJobs()
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

// GetDurationAsString returns a pretty printed duration
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
		duration += ".."
	}
	return
}

func (cmd *jobsStatusCommand) printCsv() error {
	fmt.Printf("folder,name,status,duration,starttime,endtime\n")
	for _, job := range cmd.reply.Statuses {
		fmt.Printf("%s,%s,%s,%f,%s,%s\n", job.Folder, job.Name, job.Status, GetDuration(job), job.StartTime, job.EndTime)
	}
	return nil
}

func (cmd *jobsStatusCommand) PrettyPrint() error {
	if cmd.Csv {
		return cmd.printCsv()
	}
	if cmd.Verbose {
		fmt.Printf("%d/%d jobs displayed\n", cmd.reply.Returned, cmd.reply.Total)
		fmt.Printf("%-40.40s %5.5s %-14.14s %8.8s %14.14s %17.17s %5.5s %12.12s %21.21s %8.8s %4.4s\n",
			"Folder/Name", "Held", "JobId", "Order", "Status", "Host", "Del?", "Start time", "Description", "Duration", "Runs")
		fmt.Printf("%s\n", strings.Repeat("-", 158))
		for _, job := range cmd.reply.Statuses {
			fmt.Printf("%-40.40s %5.5s %-14.14s %8.8s %14.14s %17.17s %5.5s %12.12s %21.21s %8.8s % 4d\n",
				job.Folder+"/"+job.Name,
				strconv.FormatBool(job.Held),
				job.JobId, job.OrderDate, job.Status, job.Host, strconv.FormatBool(job.Deleted), job.StartTime, job.Description, GetDurationAsString(job), job.NumberOfRuns)
		}
	} else {
		fmt.Printf("%-15.15s %18.18s %8.8s %16.16s %8.8s\n",
			"Name", "JobId", "Status", "Host", "Duration")
		fmt.Printf("---------------------------------------------------------------------\n")
		for _, job := range cmd.reply.Statuses {
			fmt.Printf("%-15.15s %18.18s %8.8s %18.18s %6.6s\n",
				job.Name,
				job.JobId, job.Status, job.Host, GetDurationAsString(job))
		}

	}
	return nil
}
func init() {
	commands.AddCommand("lj", "list the jobs", "List jobs matching the filtering criteria", &jobsStatusCommand{})
}
