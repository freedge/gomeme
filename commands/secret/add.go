// Package secret contains a few commands to handle secrets
package secret

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func readSecretFromFile(filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContents, nil
}

type add struct {
	Name  string `short:"n" long:"name" description:"name" required:"true"`
	File  string `short:"f" long:"file" description:"secret is the content of that file. Use - to read from standard input" required:"true"`
	reply types.Message
}

func (cmd *add) Execute([]string) (err error) {
	if err := commands.RequiresAnnotation(); err != nil {
		return err
	}
	var content []byte
	if content, err = readSecretFromFile(cmd.File); err != nil {
		return err
	}
	err = client.Call("POST", "/config/secret", types.SecretAddQuery{Name: cmd.Name, Value: string(content)}, map[string]string{}, &cmd.reply)
	return
}

func (cmd *add) Data() interface{} {
	return cmd.reply
}

func (cmd *add) PrettyPrint() error {
	fmt.Println(cmd.reply.Message)
	return nil
}

func init() {
	commands.AddCommand("secret.add", "add secret", "Add a named secret", &add{})
}
