package curl

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/commands"
)

// the powershell version of curl
type invokeRestMethod struct {
	out string
}

func (cmd *invokeRestMethod) Data() interface{} {
	return &cmd.out
}

func (cmd *invokeRestMethod) Execute([]string) error {
	var kflag string
	// powershell does not have the equivalent of capath sadly
	if commands.Opts.Capath != "" {
		kflag = " -SkipCertificateCheck"
	}
	headers := []string{"Accept='application/json'"}
	if theToken, found := commands.Tokens.Endpoint[commands.Opts.Endpoint]; found {
		headers = append(headers, fmt.Sprintf("Authorization='Bearer %s'", theToken.Token.Token))
	}

	cmd.out = fmt.Sprintf("Invoke-RestMethod%s -Headers @{%s}  -Uri %s",
		kflag, strings.Join(headers, "; "), commands.Opts.Endpoint)
	return nil
}
func (cmd *invokeRestMethod) PrettyPrint() error {
	fmt.Println(cmd.out)
	return nil
}

func init() {
	commands.AddCommand("ps", "PowerShell 7 curl command", "The Invoke-RestMethod command to type to access the API under PowerShell 7", &invokeRestMethod{})
}
