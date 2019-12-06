package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type put struct {
	Filename string `long:"filename" short:"f" description:"json file corresponding to the job definition folder" required:"true"`
	Dry      bool   `long:"dry" short:"d" description:"call the build service, not the deploy one" `
	Ctm      string `long:"ctm" short:"c" description:"controlm server, to prevent mistakes" required:"true"`
	reply    types.DeployReply
}

func (cmd *put) Data() interface{} {
	return cmd.reply
}

// NewfileUploadRequest creates a new file upload http request with optional extra params
func NewfileUploadRequest(filename string, filecontents []byte, paramName string, writer *multipart.Writer) error {
	part, err := writer.CreateFormFile(paramName, filename)
	if err != nil {
		return nil
	}
	_, err = part.Write(filecontents)
	return err
}

type linter func([]byte) error

func getLinter(ctm string) linter {
	return func(content []byte) error {
		var file types.DeployPutFormat
		if err := json.Unmarshal(content, &file); err != nil {
			return fmt.Errorf("could not parse that file: %w", err)
		}
		for _, folder := range file {
			if folder.ControlmServer != ctm && folder.TargetCTM != ctm {
				return fmt.Errorf("controlm server does not match")
			}
			// we make an exception for the workbench...
			// TODO: either remove the sitestandard to be able to deploy the job on workbench,
			// or find a way to define the sitestandard on workbench as well (does not
			// seem possible through the API today)
			// TODO: check what we do about connection profiles
			if ctm != "workbench" && folder.SiteStandard == "" && folder.TargetCTM == "" {
				return fmt.Errorf("site standard not specified")
			}
		}
		return nil
	}
}

func newfileUploadRequestFromFile(path string, paramName string, writer *multipart.Writer, l linter) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	if err = l(fileContents); err != nil {
		return err
	}

	return NewfileUploadRequest(fi.Name(), fileContents, paramName, writer)
}

func (cmd *put) Execute([]string) (err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if err = newfileUploadRequestFromFile(cmd.Filename, "definitionsFile", writer, getLinter(cmd.Ctm)); err != nil {
		return
	}
	writer.Close()
	client.ContentType = writer.FormDataContentType()
	service := "/deploy"
	if cmd.Dry {
		service = "/build"
	} else {
		if err := commands.RequiresAnnotation(); err != nil {
			return err
		}
	}
	err = client.Call("POST", service, body, map[string]string{}, &cmd.reply)
	return
}

func (cmd *put) PrettyPrint() error {
	for _, deploy := range cmd.reply {
		fmt.Printf("%d jobs successfully deployed (%#v)\n", deploy.SuccessfulJobsCount, deploy)
		for _, oneerr := range deploy.Errors {
			fmt.Printf("Got error : %s\n", strings.Join(oneerr.Lines, ", "))
		}
	}
	return nil
}

func init() {
	commands.AddCommand("deploy.put", "upload", "Upload all the jobs", &put{})
}
