// The watchclock command updates Object Lock retention periods for
// objects in an S3 bucket.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mechfish/subcommander"
	"github.com/mechfish/watchclock"
)

// Transform the subcommander Config and pass it to the command.
func runCommand(cmd func(*watchclock.Config, []string) error) func(subcommander.Config, []string) error {
	return func(sConf subcommander.Config, args []string) error {
		if config, ok := sConf.(*watchclock.Config); ok {
			if err := config.Validate(); err != nil {
				return err
			}
			defer config.Cleanup()
			return cmd(config, args)
		}
		panic("Configuration error: Could not convert CLI config")
	}
}

var commands = &subcommander.CommandSet{
	Name:               "watchclock",
	DefaultCommandName: "renew",
	Commands: []subcommander.Command{
		{
			Name:            "renew",
			Description:     "Update S3 Object Locks that are nearly expired",
			NumArgsRequired: 1,
			Run:             runCommand(watchclock.Renew),
		},
	},
}

func main() {
	err := commands.Execute(&watchclock.Config{})
	if err == nil {
		return
	}
	var helpE *subcommander.NeededHelpError
	if errors.As(err, &helpE) {
		os.Exit(2)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "%s\n", err.Error())
	os.Exit(1)
}
