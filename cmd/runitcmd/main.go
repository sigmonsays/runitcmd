package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/sigmonsays/runitcmd/runit"

	gologging "github.com/sigmonsays/go-logging"
)

type Application struct {
	*cli.App
	Runit *runit.Runit
}

func main() {
	c := cli.NewApp()
	c.Name = "runitcmd"
	c.Version = "0.0.1"
	app := &Application{
		App: c,
	}

	rcfg := runit.DefaultRunitConfig()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "level, l",
			Value: "WARN",
			Usage: "change log level",
		},
		cli.StringFlag{
			Name:  "service-dir",
			Usage: "change service dir",
		},
		cli.StringFlag{
			Name:  "active-dir",
			Usage: "change active service dir",
		},
	}

	app.Before = func(c *cli.Context) error {
		gologging.SetLogLevel(c.String("level"))

		service_dir := c.String("service-dir")
		active_dir := c.String("active-dir")
		if service_dir != "" {
			rcfg.ServiceDir = service_dir
		}
		if active_dir != "" {
			rcfg.ActiveDir = active_dir
		}
		app.Runit = runit.NewRunit(rcfg)

		return nil
	}

	app.Commands = []cli.Command{
		initList(app),
	}

	app.Run(os.Args)
}
