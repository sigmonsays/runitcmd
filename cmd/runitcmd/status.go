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
		fmt.Printf("\n")
	}
}
