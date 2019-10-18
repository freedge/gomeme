package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type LoginCommand struct {
	Username string
}

const (
	PASSWORD = "GOMEME_PASSWORD"
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
	if Endpoint == "" || password == "" || cmd.Username == "" {
		err = fmt.Errorf("Endpoint (%s), password (set=%t), username (%s) must be set", Endpoint, password != "", cmd.Username)
		return
	}

	query := SessionLoginQuery{Username: cmd.Username, Password: password}
	jsonquery, err := json.Marshal(query)
	if err != nil {
		return
	}
	resp, err := http.Post(Endpoint+"/session/login", "application/json", bytes.NewBuffer(jsonquery))
	if err != nil {
		return
	}
	var token Token
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return
	}
	i = token.Token
	ioutil.WriteFile(".token", []byte(token.Token), 0600)

	return

}

func (cmd *LoginCommand) PrettyPrint(flags *flag.FlagSet, i interface{}) error {
	fmt.Println("token", i)
	return nil
}

func init() {
	Register("login", &LoginCommand{})
}
