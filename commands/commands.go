// Package commands is the base of all commands. All commands
// should be in a sub package, implement the Command interface,
// and register itself through the init function using the Register function
package commands

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"
)

// DefaultOpts are the common options for every command
type DefaultOpts struct {
	Dump       bool `long:"dump" description:"outputs as go structure"`
	JSONNeeded bool `long:"json" description:"outputs as json"`
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
	Data() interface{}  // Return the data after the command ran succesfully
	PrettyPrint() error // Pretty print the output of the command. It is given the data as returned by the Run method
}

// AddIfNotEmpty is a convenient method to add in a map if value is not empty
func AddIfNotEmpty(args map[string]string, key, value string) {
	if value != "" {
		args[key] = value
	}
}

// Usage prints the usage of this function using the registered command names
func Usage() {
	// 	var sb strings.Builder
	// 	sb.WriteString("Usage: ")
	// 	sb.WriteString(os.Args[0])
	// 	sb.WriteString(" [command] -h\n Supported commands:\n")

	// 	var keys []string
	// 	for key := range Commands {
	// 		keys = append(keys, key)
	// 	}
	// 	sort.Strings(keys)
	// 	for _, key := range keys {
	// 		fmt.Fprintf(&sb, "\t%s\n", key)
	// 	}

	// 	fmt.Fprintf(&sb, `
	// export %s environment variable to the Control-M endpoint (eg: https://foobar:8443/automation-api)
	// export %s=true environment variable to skip host verification
	// `, envEndpoint, envInsecure)

	// fmt.Println(sb.String())
}

const (
	envEndpoint = "GOMEME_ENDPOINT" // environment variable for the URL to target
	envInsecure = "GOMEME_INSECURE" // environment variable when the certifacte of the endpoint is not properly set up
)

var (
	// Endpoint URL to target
	Endpoint string

	// Insecure is set to skip SSL verify
	Insecure bool

	// TheToken to use to connect to the endpoint
	TheToken string
)

// init sets us endpoint and token
func init() {
	Endpoint = os.Getenv(envEndpoint)
	Insecure, _ = strconv.ParseBool(os.Getenv(envInsecure))

	if s, err := ioutil.ReadFile(".token"); err == nil {
		TheToken = string(s)
	}
}
