// Package commands is the base of all commands. All commands
// should be in a sub package, implement the Command interface,
// and register itself through the init function using the Register function
package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Command is implemented by all our commands
type Command interface {
	Prepare(flags *flag.FlagSet)        // Prepare registers its own set of flags
	Run() (interface{}, error)          // Run the command, return an object that can be later dump as json
	PrettyPrint(data interface{}) error // Pretty print the output of the command. It is given the data as returned by the Run method
}

// Each and every command must register itself into this map through a call to Register in the init function
var Commands map[string]Command

// Register a new command, to be called from the init function
func Register(name string, cmd Command) {
	if Commands == nil {
		Commands = make(map[string]Command, 0)
	}
	Commands[name] = cmd
}

// Usage prints the usage of this function using the registered command names
func Usage() {
	s := "Usage: " + os.Args[0] + " ["
	for key := range Commands {
		s = s + " " + key
	}
	s = s + " ]."
	if Endpoint == "" {
		s += " \nSet " + ENDPOINT + " environment variable to the Control-M endpoint (eg: https://foobar:8443/automation-api)\n"
		s += "Set " + INSECURE + " environment variable to skip host verification\n"
	}
	fmt.Println(s)
}

const (
	ENDPOINT = "GOMEME_ENDPOINT" // environment variable for the URL to target
	INSECURE = "GOMEME_INSECURE" // environment variable when the certifacte of the endpoint is not properly set up
)

var Endpoint string
var Insecure bool
var TheToken string // The token to use to connect to the endpoint

// init sets us endpoint and token
func init() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case ENDPOINT:
			Endpoint = pair[1]
		case INSECURE:
			Insecure, _ = strconv.ParseBool(pair[1])
		}
	}

	s, err := ioutil.ReadFile(".token")
	if err == nil {
		TheToken = string(s)
	}
}
