package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/freedge/gomeme/commands"
	_ "github.com/freedge/gomeme/commands/curl"
	_ "github.com/freedge/gomeme/commands/job"
	_ "github.com/freedge/gomeme/commands/login"
	_ "github.com/freedge/gomeme/commands/qr"
)

var (
	dumpNeeded = false // dump the go object
	jsonNeeded = false // dump as json
)

func main() {
	if len(os.Args) < 2 {
		commands.Usage()
		os.Exit(-1)
	}

	// find the proper command and delegate most of the actions to it
	command, found := commands.Commands[os.Args[1]]
	if !found {
		commands.Usage()
		os.Exit(-1)
	}
	flagset := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	flagset.BoolVar(&dumpNeeded, "dump", false, "outputs as go structure")
	flagset.BoolVar(&jsonNeeded, "json", false, "outputs as json")
	command.Prepare(flagset)

	if err := flagset.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	data, err := command.Run()
	if err != nil {
		fmt.Printf("command exited in error: %s\n", err.Error())
		os.Exit(-1)
	}

	switch {
	case jsonNeeded:
		bytes, _ := json.MarshalIndent(data, "", "  ")
		fmt.Printf("%s\n", string(bytes))
	case dumpNeeded:
		fmt.Printf("%#v", data)
	default:
		err = command.PrettyPrint(data)
	}
	if err != nil {
		fmt.Printf("Printing exited in error: %s\n", err.Error())
		os.Exit(-1)
	}
}
