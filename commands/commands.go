// Package commands is the base of all commands. All commands
// should be in a sub package, implement the Command interface,
// and register itself through the init function using the Register function
package commands

import (
	"github.com/jessevdk/go-flags"
)

// DefaultOpts are the common options for every command
type DefaultOpts struct {
	Dump        bool   `long:"dump" description:"outputs as go structure"`
	JSONNeeded  bool   `long:"json" description:"outputs as json"`
	Capath      string `long:"capath" description:"SSL_CERT_DIR for GOMEME"  env:"GOMEME_CERT_DIR"`
	Endpoint    string `long:"endpoint" description:"endpoint" env:"GOMEME_ENDPOINT" required:"true"`
	Debug       bool   `long:"debug"`
	Subject     string `long:"subject" description:"annotation subject"`
	Description string `long:"description" description:"annotation description"`
}

// Opts is the list of default opts
var Opts DefaultOpts

// Parser defines all the commands
var Parser = flags.NewParser(&Opts, flags.Default)

// AddCommand register a new command to be parsed
func AddCommand(a, b, c string, cmd Command) {
	Parser.AddCommand(a, b, c, cmd)
}

// Command is implemented by all our commands
type Command interface {
	flags.Commander
	Data() interface{}  // Return the data after the command ran successfully
	PrettyPrint() error // Pretty print the output of the command. It is given the data as returned by the Run method
}

// AddIfNotEmpty is a convenient method to add in a map if value is not empty
func AddIfNotEmpty(args map[string]string, key, value string) {
	if value != "" {
		args[key] = value
	}
}
