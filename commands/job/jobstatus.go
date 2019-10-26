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
	application string `short:"a" long:"application"`
	limit       int    `short:"l"  long:"limit" default:"1000"`
	status      string `short:"s"  long:"status" choice:"Executing" choice:"Ended Not OK" choice:"Ended OK" description:"Only this status"`
	reply       types.JobsStatusReply
	jobname     string `short:"n" long:"job name" description:"job name"`
	jobid       string `short:"j" long:"jobid" description:"job id"`
	folder      string `short:"f" long:"folder" description:"folder"`
	verbose     bool   `short:"v" long:"verbose" description:"output more stuff"`
	host        string `short:"h" long:"host" description:"host"`
	neighbours  bool   `long:"deps" description:"browse through neighours of this job. Only jobid can be used to filter jobs"`
}

type jobsStatusCommand struct {
	jobsStatusCommonCommand
	csv bool `short:"c" long:"csv" description:"csv output"`
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
	commands.AddIfNotEmpty(args, "application", cmd.application)
	commands.AddIfNotEmpty(args, "status", cmd.status)
	commands.AddIfNotEmpty(args, "jobname", cmd.jobname)
	commands.AddIfNotEmpty(args, "folder", cmd.folder)
	commands.AddIfNotEmpty(args, "host", cmd.host)

	if cmd.jobid != "" {
		if cmd.neighbours {
			if len(args) > 0 {
				err = fmt.Errorf("only jobid should be used to filter jobs")
				return
			}
			args["neighborhood"] = "1"
			args["direction"] = "radial"
			args["depth"] = "5"
		}
		args["jobid"] = cmd.jobid
	} else {
		if cmd.neighbours {
			err = fmt.Errorf("jobid missing")
			return
		}
	}
	args["limit"] = strconv.Itoa(cmd.limit)

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
	if cmd.csv {
		return cmd.printCsv()
	}
	if cmd.verbose {
		fmt.Printf("%d/%d jobs displayed\n", cmd.reply.Returned, cmd.reply.Total)
		fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.16s %16.16s %5.5s %12.12s %12.12s %20.20s %8.8s %4.4s\n",
			"Folder/Name", "Held", "JobId", "Order", "Status", "Host", "Del?", "Start time", "End time", "Description", "Duration", "Runs")
		fmt.Printf("%s\n", strings.Repeat("-", 177))
		for _, job := range cmd.reply.Statuses {
			fmt.Printf("%-40.40s %5.5s %-20.20s %8.8s %16.146s %16.16s %5.5s %12.12s %12.12s %20.20s %8.8s % 4d\n",
				job.Folder+"/"+job.Name,
				strconv.FormatBool(job.Held),
				job.JobId, job.OrderDate, job.Status, job.Host, strconv.FormatBool(job.Deleted), job.StartTime, job.EndTime, job.Description, GetDurationAsString(job), job.NumberOfRuns)
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
	commands.AddCommand("lj", "listjobs", "list jobs", &jobsStatusCommand{})
}
