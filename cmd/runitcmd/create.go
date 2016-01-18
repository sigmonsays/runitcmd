package main

import (
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/sigmonsays/runitcmd/runit"
)

func initCreate(app *Application) cli.Command {
	description := "Create services"
	usage := "Create services"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "exec, e",
			Usage: "execute command",
		},
		cli.StringFlag{
			Name:  "log-dir, l",
			Usage: "log to directory",
		},
		cli.BoolFlag{
			Name:  "disabled, d",
			Usage: "create service but do not enable it",
		},
	}

	cmd := cli.Command{
		Name:        "create",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Create,
	}
	return cmd
}

func (app *Application) Create(c *cli.Context) {
	args := c.Args()
	name := args.First()
	exec := c.String("exec")
	disabled := c.Bool("disabled")
	log_dir := c.String("log-dir")

	if name == "" {
		log.Errorf("service name is required")
		return
	}
	if exec == "" {
		log.Errorf("execute is required")
		return
	}

	log.Tracef("Create %s", name)
	lcfg := runit.DefaultLoggingConfig()
	if log_dir == "" {
		lcfg.Directory = filepath.Join(app.Runit.ServiceDir, name)
	} else {
		lcfg.Directory = log_dir
	}

	cfg := &runit.ServiceConfig{
		Name:           name,
		Exec:           exec,
		Logging:        lcfg,
		Disabled:       disabled,
		RedirectStderr: true,
	}

	err := app.Runit.Create(cfg)
	if err != nil {
		log.Warnf("Create %s: %s", name, err)
	}

}
