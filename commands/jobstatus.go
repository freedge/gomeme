package commands

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type JobsStatusCommand struct {
	application string
	limit       int
	status      string
	reply       JobsStatusReply
}

func (cmd *JobsStatusCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.application, "application", "", "Jobs for this application")
	flags.IntVar(&cmd.limit, "limit", 1000, "Limit to how many jobs")
	flags.StringVar(&cmd.status, "status", "", "Only this status")
}
func (cmd *JobsStatusCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	i = nil
	var bearer = "Bearer " + TheToken
	req, err := http.NewRequest("GET", Endpoint+"/run/jobs/status", nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	if cmd.application != "" {
		q.Add("application", cmd.application)
	}
	if cmd.status != "" {
		q.Add("status", cmd.status)
	}
	q.Add("limit", strconv.Itoa(cmd.limit))

	req.URL.RawQuery = q.Encode()

	// Send req using http Client
	client := &http.Client{}
	if Insecure {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
		}
		client.Transport = &http.Transport{
			TLSClientConfig: cfg,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &cmd.reply)
	if err != nil {
		return
	}
	i = cmd.reply.Statuses

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
	Register("lj", &JobsStatusCommand{})
}
