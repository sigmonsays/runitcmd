package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func initStatus(app *Application) cli.Command {
	description := "get service status"
	usage := "get service status"

	flags := []cli.Flag{}

	cmd := cli.Command{
		Name:        "status",
		Aliases:     []string{"st"},
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Status,
	}
	return cmd
}

func (app *Application) Status(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		st, err := service.Status()
		if err != nil {
			log.Warnf("%s: %s", service.Name, err)
			continue
		}

		fmt.Printf("name       %s\n", service.Name)
		fmt.Printf("running    %v\n", st.Running)
		fmt.Printf("pid        %d\n", st.Pid)
		fmt.Printf("seconds    %d (%s)\n", st.Seconds, formatTime(st.Seconds))
		fmt.Printf("enabled    %v\n", st.Enabled)
		fmt.Printf("activated  %v\n", st.Activated)
		if cfg, err := service.Config(); err == nil {
			fmt.Printf("exec       %s\n", cfg.Exec)
			fmt.Printf("\n")

			l := cfg.Logging
			if l != nil {
				fmt.Printf("logging:\n")
				fmt.Printf("directory       %s\n", l.Directory)
				fmt.Printf("max size        %d bytes\n", l.Size)
				fmt.Printf("number          %d\n", l.Num)
				fmt.Printf("timeout         %d seconds\n", l.Timeout)
				fmt.Printf("minimum         %d\n", l.Min)
				fmt.Printf("\n")
			}
		} else {
			fmt.Printf("logging configuration unavailable\n")
		}
	}
}
