package commands

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OrderJob struct {
	Hold   bool
	Ctm    string
	Folder string
	Jobs   string
	data   []byte
}

func (cmd *OrderJob) Prepare(flags *flag.FlagSet) {
	flags.BoolVar(&cmd.Hold, "hold", true, "Hold the job after submission")
	flags.StringVar(&cmd.Ctm, "ctm", "", "ctm")
	flags.StringVar(&cmd.Folder, "folder", "", "Folder")
	flags.StringVar(&cmd.Jobs, "jobs", "", "jobs")
}
func (cmd *OrderJob) Run(flags *flag.FlagSet) (i interface{}, err error) {
	i = nil
	fmt.Println("untested !")
	if TheToken == "" {
		err = fmt.Errorf("no token found. Please login first")
		return
	}
	if !cmd.Hold {
		err = fmt.Errorf("Currently only support holding jobs")
		return
	}

	if cmd.Jobs == "" || cmd.Ctm == "" || cmd.Folder == "" {
		err = fmt.Errorf("parameters missing")
		return
	}

	var bearer = "Bearer " + TheToken
	query := OrderQuery{Jobs: cmd.Jobs, Ctm: cmd.Ctm, Folder: cmd.Folder, Hold: cmd.Hold}
	jsonquery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", Endpoint+"/run/order", bytes.NewBuffer(jsonquery))

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

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
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	cmd.data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	i = cmd.data

	return
}

func (cmd *OrderJob) PrettyPrint(flags *flag.FlagSet, data interface{}) error {
	fmt.Println(data)
	return nil
}

func init() {
	Register("joborder", &OrderJob{})
}
