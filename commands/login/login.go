// Package login defines the login command that retrieves a token
// and save it into a file
package login

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type loginCommand struct {
	user  string
	token types.Token
}

const (
	envPassword = "GOMEME_PASSWORD" // environment variable for your password, only used by the login command
)

func (cmd *loginCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&(cmd.user), "user", "", "Username to use")
}

func (cmd *loginCommand) Run() (i interface{}, err error) {
	if commands.Endpoint == "" || cmd.user == "" {
		err = fmt.Errorf("Endpoint (%s) and username (%s) must be set", commands.Endpoint, cmd.user)
		return
	}

	// either get the password from the environment
	var password string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case envPassword:
			password = pair[1]
		}
	}

	// either get it on the terminal
	if password == "" {
		if !terminal.IsTerminal(0) {
			err = fmt.Errorf("run from a terminal or provide password through %s environment variable", envPassword)
			return
		}
		fmt.Printf("Enter password for user %s:\n", cmd.user)
		var bytes []byte
		bytes, err = terminal.ReadPassword(0 /* stdin */)
		if err != nil {
			return
		}
		password = string(bytes)
	}

	query := types.SessionLoginQuery{Username: cmd.user, Password: password}

	err = client.Call("POST", "/session/login", query, map[string]string{}, &cmd.token)
	if err != nil {
		return
	}

	i = cmd.token
	err = ioutil.WriteFile(".token", []byte(cmd.token.Token), 0600)
	return
}

func (cmd *loginCommand) PrettyPrint(i interface{}) error {
	fmt.Printf("Logged in. Server version %s\n", cmd.token.Version)
	return nil
}

func init() {
	commands.Register("login", &loginCommand{})
}
