package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/freedge/gomeme/commands"
)

// Call the specific url under endpoint, with the proper query and parameters
func Call(method, url string, query interface{}, params map[string]string, out interface{}) (err error) {
	var bytebuffer io.Reader = nil
	if query != nil {
		jsonquery, err := json.Marshal(query)
		if err != nil {
			return err
		}
		bytebuffer = bytes.NewBuffer(jsonquery)
	}
	req, err := http.NewRequest(method, commands.Endpoint+url, bytebuffer)

	if commands.TheToken != "" {
		var bearer = "Bearer " + commands.TheToken

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)
	}

	q := req.URL.Query()
	for name, value := range params {
		q.Add(name, value)
	}
	req.URL.RawQuery = q.Encode()

	// Send req using http Client
	client := &http.Client{}
	if commands.Insecure {
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
	switch resp.StatusCode {
	case 404:
		err = fmt.Errorf("client: got an error accessing %v", req.URL)
		return

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch out.(type) {
	case *string:
		*(out.(*string)) = string(body)
	default:
		err = json.Unmarshal(body, out)
	}

	return
}
