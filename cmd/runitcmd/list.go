package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func initList(app *Application) cli.Command {
	description := "list services"
	usage := "list services"

	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "list all services",
		},
	}

	cmd := cli.Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.List,
	}
	return cmd
}

func (app *Application) List(c *cli.Context) {
	show_all := c.Bool("all")

	log.Tracef("list all:%v", show_all)

	services, err := app.Runit.ListServices()
	if err != nil {
		log.Warnf("ListServices: %s", err)
		return
	}

	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, service := range services {
		if service.Exists() == false {
			continue
		}
		if show_all == false && service.Enabled() == false {
			continue
		}
		st, err := service.Status()
		if err != nil {
			log.Warnf("Status: %s: %s", service.Name, err)
			continue
		}
		fmt.Fprintf(tw, "%s \t running:%v \t pid:%d \t sec:%d\n", service.Name, st.Running, st.Pid, st.Seconds)
	}
	tw.Flush()

}
