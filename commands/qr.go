package commands

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type QRCommand struct {
	Name string
	Ctm  string
	qrs  []QR
}

func (cmd *QRCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&cmd.Name, "name", "", "resource name")
	flags.StringVar(&cmd.Ctm, "ctm", "", "ctm")
}
func (cmd *QRCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	i = nil
	if TheToken == "" {
		err = fmt.Errorf("no token found. Please login first")
		return
	}

	var bearer = "Bearer " + TheToken
	req, err := http.NewRequest("GET", Endpoint+"/run/resources", nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	if cmd.Name != "" {
		q.Add("name", cmd.Name)
	}
	if cmd.Ctm != "" {
		q.Add("ctm", cmd.Ctm)
	}

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

	err = json.Unmarshal(body, &cmd.qrs)
	if err != nil {
		return
	}
	i = &cmd.qrs

	return
}

func (cmd *QRCommand) PrettyPrint(f *flag.FlagSet, i interface{}) error {
	fmt.Println("QR                       Ctm        Available      Max")
	fmt.Println("======================================================")
	for _, qr := range cmd.qrs {
		fmt.Printf("%-24.24s %-8.8s %11.11s %8d\n", qr.Name, qr.Ctm, qr.Available, qr.Max)
	}
	return nil
}

func init() {
	Register("qr", &QRCommand{})
}
