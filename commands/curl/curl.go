/*
Package curl defines the curl command

gomeme does not provide everything, so in last resort one can still
call the API by hand through curl. The gomeme curl command outputs
the right curl command line to use
*/
package curl

import (
	"fmt"

	"github.com/freedge/gomeme/commands"
)

type curl struct {
	out string
}

func (cmd *curl) Data() interface{} {
	return &cmd.out
}

func (cmd *curl) Execute([]string) error {
	var kflag string
	if commands.Insecure {
		kflag = "-k "
	}
	cmd.out = fmt.Sprintf("curl %s-H 'Accept: application/json' -H 'Authorization: Bearer %s' %s",
		kflag, commands.TheToken, commands.Endpoint)
	return nil
}
func (cmd *curl) PrettyPrint() error {
	fmt.Println(cmd.out)
	return nil
}

func init() {
	commands.AddCommand("curl", "", "", &curl{})
}
