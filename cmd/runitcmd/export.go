package main

import (
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func initExport(app *Application) *cli.Command {
	description := "export service"
	usage := "export service"

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "forgiving",
			Aliases: []string{"f"},
			Usage:   "be forgiving and try to parse the runit files",
		},
	}

	cmd := &cli.Command{
		Name:        "export",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Export,
	}
	return cmd
}

func (app *Application) Export(c *cli.Context) error {
	forgiving := c.Bool("forgiving")

	for _, service := range app.MatchingServices(c) {

		cfg, err := app.Runit.Export(service.Name)
		if forgiving == false && err != nil {
			log.Warnf("export %s: %s", service.Name, err)
			continue
		}
		if forgiving && err != nil {
			cfg, err = app.Runit.LoadFromDisk(service.Name)
			if err != nil {
				log.Warnf("load service config %s: %s", service.Name, err)
				continue
			}
		}

		destfile := filepath.Join("./", service.Name+".yaml")
		err = cfg.SaveFile(destfile)
		if err != nil {
			log.Warnf("export %s: %s", service.Name, err)
			continue
		}
	}
	return nil
}
