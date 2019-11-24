package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/commands/deploy"
	"github.com/freedge/gomeme/types"

	"github.com/tidwall/sjson"
)

type modify struct {
	Jobid  string `short:"j" long:"jobid" description:"job id" required:"true"`
	Name   string `short:"n" long:"name" description:"job name, only there to avoid mistakes" required:"true"`
	jobdef string
	reply  string
	job    types.JobGetReply
}

func (cmd *modify) Data() interface{} {
	return cmd.job
}

func (cmd *modify) Execute(args []string) (err error) {
	// starts by getting the job definition
	err = client.Call("GET", "/run/job/"+cmd.Jobid+"/get", nil, map[string]string{}, &cmd.jobdef)
	if err != nil {
		return err
	}

	// we're doing the parsing ourselves because we want to patch the output, not regenerate a new definition from scratch
	err = json.Unmarshal([]byte(cmd.jobdef), &cmd.job)
	if err != nil {
		return err
	}
	if len(cmd.job) != 1 {
		return fmt.Errorf("Did not find 1 job definition")
	}

	// do some magic modification
	for name, def := range cmd.job {
		if name != cmd.Name {
			return fmt.Errorf("Job name is %s and not %s", name, cmd.Name)
		}
		fmt.Println("Replacing:", def.Arguments, "with", args)
		cmd.jobdef, err = sjson.Set(cmd.jobdef, fmt.Sprintf("%s.Arguments", name), args)
		if err != nil {
			return
		}
	}

	// make subject mandatory
	if commands.Opts.Subject == "" {
		fmt.Println("Skipping because no subject provided")
		return nil
	}

	// now upload the modified job definition
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if err = deploy.NewfileUploadRequest("foo.json", []byte(cmd.jobdef), "jobDefinitionsFile", writer); err != nil {
		return
	}
	writer.Close()
	client.ContentType = writer.FormDataContentType()

	err = client.Call("POST", "/run/job/"+cmd.Jobid+"/modify", body, map[string]string{}, &cmd.reply)

	return err
}

func (cmd *modify) PrettyPrint() error {
	fmt.Println(cmd.reply)
	return nil
}

func init() {
	commands.AddCommand("job.modify", "modify", "modify parameters of a held job", &modify{})
}
