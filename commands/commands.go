// Package commands is the base of all commands. All commands
// should be in a sub package, implement the Command interface,
// and register itself through the init function using the Register function
package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Command is implemented by all our commands
type Command interface {
	Prepare(flags *flag.FlagSet)        // Prepare registers its own set of flags
	Run() (interface{}, error)          // Run the command, return an object that can be later dump as json
	PrettyPrint(data interface{}) error // Pretty print the output of the command. It is given the data as returned by the Run method
}

// Commands contains Each and every command, must registered through a call to Register in the init function
var Commands map[string]Command

// AddIfNotEmpty is a convenient method to add in a map if value is not empty
func AddIfNotEmpty(args map[string]string, key, value string) {
	if value != "" {
		args[key] = value
	}
}

// Register a new command, to be called from the init function
func Register(name string, cmd Command) {
	if Commands == nil {
		Commands = make(map[string]Command, 0)
	}
	Commands[name] = cmd
}

// Usage prints the usage of this function using the registered command names
func Usage() {
	var sb strings.Builder
	sb.WriteString("Usage: ")
	sb.WriteString(os.Args[0])
	sb.WriteString(" [command] -h\n Supported commands:\n")

	var keys []string
	for key := range Commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&sb, "\t%s\n", key)
	}

	fmt.Fprintf(&sb, `
export %s environment variable to the Control-M endpoint (eg: https://foobar:8443/automation-api)
export %s=true environment variable to skip host verification
`, envEndpoint, envInsecure)

	fmt.Println(sb.String())
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
