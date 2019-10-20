// Package login defines the login command that retrieves a token
// and save it into a file
package login

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type LoginCommand struct {
	Username string
}

const (
	PASSWORD = "GOMEME_PASSWORD" // environment variable for your password, only used by the login command
)

func (cmd *LoginCommand) Prepare(flags *flag.FlagSet) {
	flags.StringVar(&(cmd.Username), "username", "", "Username to use")
}

func (cmd *LoginCommand) Run(flags *flag.FlagSet) (i interface{}, err error) {
	var password string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case PASSWORD:
			password = pair[1]
		}
	}
	if commands.Endpoint == "" || password == "" || cmd.Username == "" {
		err = fmt.Errorf("Endpoint (%s), password (set=%t), username (%s) must be set", commands.Endpoint, password != "", cmd.Username)
		return
	}

	query := types.SessionLoginQuery{Username: cmd.Username, Password: password}
	var token types.Token

	err = client.Call("POST", "/session/login", query, map[string]string{}, &token)
	if err != nil {
		return
	}

	i = token.Token
	err = ioutil.WriteFile(".token", []byte(token.Token), 0600)
	return
}

func (cmd *LoginCommand) PrettyPrint(flags *flag.FlagSet, i interface{}) error {
	fmt.Println("token", i)
	return nil
}

func init() {
	commands.Register("login", &LoginCommand{})
}
