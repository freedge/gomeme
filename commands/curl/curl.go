/*
Package curl defines the curl command

gomeme does not provide everything, so in last resort one can still
call the API by hand through curl. The gomeme curl command outputs
the right curl command line to use
*/
package curl

import (
	"flag"
	"fmt"

	"github.com/freedge/gomeme/commands"
)

type curl struct {
	out string
}

func (cmd *curl) Prepare(flags *flag.FlagSet) {}
func (cmd *curl) Run(flags *flag.FlagSet) (interface{}, error) {
	kflag := ""
	if commands.Insecure {
		kflag = "-k "
	}
	cmd.out = fmt.Sprintf("curl %s-H 'Accept: application/json' -H 'Authorization: Bearer %s' %s",
		kflag, commands.TheToken, commands.Endpoint)
	return cmd.out, nil
}
func (cmd *curl) PrettyPrint(flags *flag.FlagSet, data interface{}) error {
	fmt.Println(cmd.out)
	return nil
}

func init() {
	commands.Register("curl", &curl{})
}
