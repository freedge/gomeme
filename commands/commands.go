package commands

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Command interface {
	Prepare(flags *flag.FlagSet)
	Run(flags *flag.FlagSet) (interface{}, error)
	PrettyPrint(flags *flag.FlagSet, data interface{}) error
}

var Commands map[string]Command

func Register(name string, cmd Command) {
	if Commands == nil {
		Commands = make(map[string]Command, 0)
	}
	Commands[name] = cmd
}

func Usage() {
	s := "Usage: " + os.Args[0] + " ["
	for key, _ := range Commands {
		s = s + " " + key
	}
	s = s + " ]. \n"
	s += "Set " + ENDPOINT + " environment variable to the Control-M endpoint (eg: https://foobar:8443/automation-api)\n"
	s += "Set " + INSECURE + " environment variable to skip host verification\n"
	fmt.Println(s)
}

const (
	ENDPOINT = "GOMEME_ENDPOINT"
	INSECURE = "GOMEME_INSECURE"
)

var Endpoint string
var Insecure bool
var TheToken string

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
	if Insecure {
		cfg := &tls.Config{
			InsecureSkipVerify: true,
		}
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: cfg,
		}
	}
	s, err := ioutil.ReadFile(".token")
	if err == nil {
		TheToken = string(s)
	}
}
