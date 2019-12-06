// Package secret contains a few commands to handle secrets
package secret

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
)

type get struct {
	reply []string
}

func (cmd *get) Execute([]string) (err error) {
	err = client.Call("GET", "/config/secrets", nil, map[string]string{}, &cmd.reply)
	return
}

func (cmd *get) Data() interface{} {
	return cmd.reply
}

func (cmd *get) PrettyPrint() error {
	fmt.Println("Secrets:", strings.Join(cmd.reply, ", "))
	return nil
}

func init() {
	commands.AddCommand("secret.get", "get secrets", "Get all secrets", &get{})
}
