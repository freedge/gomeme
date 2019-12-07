// Package secret contains a few commands to handle secrets
package secret

import (
	"fmt"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type update struct {
	add
}

func (cmd *update) Execute([]string) (err error) {
	if err := commands.RequiresAnnotation(); err != nil {
		return err
	}
	var content []byte
	if content, err = readSecretFromFile(cmd.File); err != nil {
		return err
	}
	err = client.Call("POST", "/config/secret/"+cmd.Name, types.SecretAddQuery{Value: string(content)}, map[string]string{}, &cmd.reply)
	return
}

func (cmd *update) Data() interface{} {
	return cmd.reply
}

func (cmd *update) PrettyPrint() error {
	fmt.Println(cmd.reply.Message)
	return nil
}

func init() {
	commands.AddCommand("secret.update", "change secret", "Change a named secret", &update{})
}
