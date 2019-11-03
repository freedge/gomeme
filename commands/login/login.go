// Package login defines the login command that retrieves a token
// and save it into a file
package login

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

type loginCommand struct {
	User  string `short:"u" long:"user" description:"Username to use" env:"USER"`
	token types.Token
}

const (
	envPassword = "GOMEME_PASSWORD" // environment variable for your password, only used by the login command
)

func (cmd *loginCommand) Data() interface{} {
	return &cmd.token
}

func (cmd *loginCommand) Execute([]string) (err error) {

	// either get it on the terminal
	var found bool
	var password string
	if password, found = os.LookupEnv(envPassword); !found {
		if !terminal.IsTerminal(0) {
			err = fmt.Errorf("run from a terminal or provide password through %s environment variable", envPassword)
			return
		}
		fmt.Printf("Enter password for user %s:\n", cmd.User)
		var bytes []byte
		if bytes, err = terminal.ReadPassword(0 /* stdin */); err != nil {
			return
		}
		password = string(bytes)
	}

	query := types.SessionLoginQuery{Username: cmd.User, Password: password}

	if err = client.Call("POST", "/session/login", query, map[string]string{}, &cmd.token); err != nil {
		return
	}

	commands.Tokens.Endpoint[commands.Opts.Endpoint] = types.TokenFileToken{Token: cmd.token, Created: time.Now()}
	commands.WriteTokensFile()

	return
}

func (cmd *loginCommand) PrettyPrint() error {
	fmt.Printf("Logged in. Server version %s\n", cmd.token.Version)
	return nil
}

func init() {
	commands.AddCommand("login", "login", "Login with the specified username", &loginCommand{})
}
