package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/freedge/gomeme/commands"
	"github.com/jessevdk/go-flags"

	_ "github.com/freedge/gomeme/commands/config"
	_ "github.com/freedge/gomeme/commands/curl"
	_ "github.com/freedge/gomeme/commands/deploy"
	_ "github.com/freedge/gomeme/commands/job"
	_ "github.com/freedge/gomeme/commands/login"
	_ "github.com/freedge/gomeme/commands/qr"
	_ "github.com/freedge/gomeme/commands/secret"
)

func commandHandler(command flags.Commander, args []string) error {
	cmd := command.(commands.Command)
	var err error
	if err = cmd.Execute(args); err != nil {
		return err
	}
	data := cmd.Data()

	switch {
	case commands.Opts.JSONNeeded:
		bytes, _ := json.MarshalIndent(data, "", "  ")
		fmt.Printf("%s\n", string(bytes))
	case commands.Opts.Dump:
		fmt.Printf("%#v", data)
	default:
		err = cmd.PrettyPrint()
	}

	return err
}

func main() {
	commands.Parser.CommandHandler = commandHandler
	if _, err := commands.Parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
