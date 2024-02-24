package main

import (
	"path/filepath"
	"strings"

	"github.com/sigmonsays/runitcmd/runit"
	"github.com/urfave/cli/v2"
)

func initApply(app *Application) *cli.Command {
	description := "update service"
	usage := "update service"

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart service after changes",
		},
	}

	cmd := &cli.Command{
		Name:        "apply",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Apply,
	}
	return cmd
}

func (app *Application) Apply(c *cli.Context) error {
	args := c.Args()
	filename := args.First()
	restart := c.Bool("restart")

	if strings.Contains(filename, "/") == false {
		filename = filepath.Join(app.Runit.ServiceDir, filename, "service.yaml")
	}

	log.Tracef("apply settings from %s", filename)

	cfg := &runit.ServiceConfig{}

	err := cfg.LoadFile(filename)
	if err != nil {
		log.Warnf("LoadFile %s: %s", filename, err)
		return err
	}

	err = app.Runit.Apply(cfg)
	if err != nil {
		log.Warnf("Create %s: %s", cfg.Name, err)
		return err
	}

	if restart {
		err = app.Runit.Restart(cfg.Name)
		if err != nil {
			log.Warnf("Restart %s: %s", cfg.Name, err)
			return err
		}

	}
	return nil
}
