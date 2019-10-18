package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/freedge/gomeme/commands"
)

var DumpNeeded = false

func main() {
	if len(os.Args) < 2 {
		commands.Usage()
		os.Exit(-1)
	}

	command, found := commands.Commands[os.Args[1]]
	if !found {
		commands.Usage()
		os.Exit(-1)
	}
	flagset := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	flagset.BoolVar(&DumpNeeded, "dump", false, "outputs as go structure")
	command.Prepare(flagset)

	err := flagset.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	data, err := command.Run(flagset)
	if err != nil {
		fmt.Printf("command exited in error: %s\n", err.Error())
		os.Exit(-1)
	}

	if DumpNeeded {
		fmt.Printf("%#v", data)
	} else {
		err = command.PrettyPrint(flagset, data)
	}
	if err != nil {
		fmt.Printf("Printing exited in error: %s\n", err.Error())
		os.Exit(-1)
	}
}
