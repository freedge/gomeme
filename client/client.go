// Package client exports methods to Call the API
package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func handleError(resp *http.Response, body []byte) (formattedError error) {
	formattedError = fmt.Errorf("server replied an error %d", resp.StatusCode)

	var reply types.ErrorReply
	err := json.Unmarshal(body, &reply)
	if err != nil {
		return
	}
	if len(reply.Errors) > 0 {
		errorString := ""
		for id, msg := range reply.Errors {
			if id > 0 {
				errorString += "\n"
			}
			errorString += msg.Message + " " + msg.Item
			if msg.Line > 0 {
				errorString += fmt.Sprintf(" (%s:%d:%d)", msg.File, msg.Line, msg.Col)
			}
			errorString += " "
		}
		formattedError = fmt.Errorf("%s", errorString)
	}

	return
}

var customTransport *http.Transport

// ContentType is the content type to send when send multi part data... maybe refactor that one day
var ContentType string

// Call the specific url under endpoint, with the proper query and parameters
func Call(method, url string, query interface{}, params map[string]string, out interface{}) (err error) {
	var bytebuffer io.Reader
	isJsonInput := true
	if query != nil {
		switch query.(type) {
		case *bytes.Buffer:
			isJsonInput = false
			bytebuffer = query.(*bytes.Buffer)
		default:
			jsonquery, err := json.Marshal(query)
			if err != nil {
				return err
			}
			bytebuffer = bytes.NewBuffer(jsonquery)
		}
	}
	req, err := http.NewRequest(method, commands.Opts.Endpoint+url, bytebuffer)

	if token, found := commands.Tokens.Endpoint[commands.Opts.Endpoint]; found {
		var bearer = "Bearer " + token.Token.Token

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
	if commands.Opts.Capath != "" {
		if customTransport == nil {
			roots := x509.NewCertPool()
			fis, err := ioutil.ReadDir(commands.Opts.Capath)
			if err != nil {
				panic(err)
			}
			for _, fi := range fis {
				if data, err := ioutil.ReadFile(commands.Opts.Capath + "/" + fi.Name()); err == nil {
					roots.AppendCertsFromPEM(data)
				}
			}
			cfg := &tls.Config{
				RootCAs: roots,
			}
			customTransport = &http.Transport{
				TLSClientConfig: cfg,
			}
		}
		client.Transport = customTransport
	}
	if isJsonInput {
		req.Header.Set("Content-type", "application/json")
	} else {
		req.Header.Set("Content-type", ContentType)
	}

	if commands.Opts.Subject != "" {
		req.Header.Set("Annotation-Subject", commands.Opts.Subject)
		req.Header.Set("Annotation-Description", commands.Opts.Description)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if commands.Opts.Debug {
		fmt.Println("DEBUG", resp.StatusCode, string(body))
	}

	switch resp.StatusCode {
	case 404:
		err = fmt.Errorf("client: got an error accessing %v", req.URL)
		return
	case 400, 401, 500, 502:
		err = handleError(resp, body)
		return
	}

	switch out.(type) {
	case *string:
		*(out.(*string)) = string(body)
	default:
		err = json.Unmarshal(body, out)
		if err != nil {
			err = handleError(resp, body)
		}
	}

	return
}
