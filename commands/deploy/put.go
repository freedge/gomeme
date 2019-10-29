package deploy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type put struct {
	Filename string `long:"filename" short:"f" description:"json file corresponding to the job definition folder"`
	reply    types.DeployReply
}

func (cmd *put) Data() interface{} {
	return cmd.reply
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(path string, paramName string, writer *multipart.Writer) error {
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

	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil
	}
	_, err = part.Write(fileContents)
	return err
}

func (cmd *put) Execute([]string) (err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if err = newfileUploadRequest(cmd.Filename, "definitionsFile", writer); err != nil {
		return
	}
	writer.Close()
	client.ContentType = writer.FormDataContentType()
	err = client.Call("POST", "/deploy", body, map[string]string{}, &cmd.reply)
	return
}

func (cmd *put) PrettyPrint() error {
	for _, deploy := range cmd.reply {
		fmt.Printf("%d jobs successfully deployed: %v\n", deploy.SuccessfulJobsCount, deploy)
	}
	return nil
}

func init() {
	commands.AddCommand("deploy.put", "upload", "Upload all the jobs", &put{})
}
