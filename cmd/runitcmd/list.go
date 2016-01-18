package main

import (
	"fmt"
	"os"
	"strings"
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

var timeLegend = []string{"w", "d", "h", "m", "s"}
var timeUnits = []int64{86400 * 7, 86400, 3600, 60, 1}

func formatTime(seconds int64) string {
	r := seconds
	ret := []string{}
	for n := 0; n < len(timeLegend); n++ {
		u := timeUnits[n]
		if r >= u {
			ret = append(ret, fmt.Sprintf("%d%s", r/u, timeLegend[n]))
			r = r - u*int64(r/u)
		}
		if r == 0 {
			break
		}
	}
	if len(ret) == 0 {
		return "0s"
	}
	return strings.Join(ret, "")
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
		pid := "-"
		state := "UP"
		name := service.Name
		if st.Running == false {
			state = "DOWN"
		} else {
			pid = fmt.Sprintf("%d", st.Pid)
		}
		fmt.Fprintf(tw, "%s\t %s \t %s \t %s\n", name, state, pid, formatTime(st.Seconds))
	}
	tw.Flush()

}
