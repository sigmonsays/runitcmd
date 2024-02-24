package main

import (
	"fmt"
	"path/filepath"

	"github.com/sigmonsays/runitcmd/runit"
	"github.com/urfave/cli/v2"
)

func initCreate(app *Application) *cli.Command {
	description := "Create services"
	usage := "Create services"

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "exec",
			Aliases: []string{"e"},
			Usage:   "execute command",
		},
		&cli.StringFlag{
			Name:    "log-dir",
			Aliases: []string{"l"},
			Usage:   "log to directory",
		},
		&cli.BoolFlag{
			Name:    "disabled",
			Aliases: []string{"d"},
			Usage:   "create service but do not enable it",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "force update the service if it already exists",
		},
		&cli.BoolFlag{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart the service after creation if it already exists",
		},
	}

	cmd := &cli.Command{
		Name:        "create",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Create,
	}
	return cmd
}

func (app *Application) Create(c *cli.Context) error {
	args := c.Args()
	name := args.First()
	exec := c.String("exec")
	force := c.Bool("force")
	restart := c.Bool("restart")
	disabled := c.Bool("disabled")
	log_dir := c.String("log-dir")

	if name == "" {
		log.Errorf("service name is required")
		return fmt.Errorf("service name is required")
	}
	if exec == "" {
		log.Errorf("execute is required")
		return fmt.Errorf("execute is required")
	}

	log.Tracef("Create %s", name)
	lcfg := runit.DefaultLoggingConfig()
	if log_dir == "" {
		lcfg.Directory = filepath.Join(app.Conf.Logging.Directory, name)
	} else {
		lcfg.Directory = log_dir
	}

	cfg := &runit.ServiceConfig{
		Name:           name,
		Exec:           exec,
		Logging:        lcfg,
		Disabled:       disabled,
		Activated:      true,
		RedirectStderr: true,
	}

	copts := &runit.CreateOptions{
		Force:   force,
		Restart: restart,
	}

	err := app.Runit.Create2(cfg, copts)
	if err != nil {
		log.Warnf("Create %s: %s", name, err)
	}

	return err
}
