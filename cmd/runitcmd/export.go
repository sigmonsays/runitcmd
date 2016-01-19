package main

import (
	"path/filepath"

	"github.com/codegangsta/cli"
)

func initExport(app *Application) cli.Command {
	description := "export service"
	usage := "export service"

	flags := []cli.Flag{}

	cmd := cli.Command{
		Name:        "export",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Export,
	}
	return cmd
}

func (app *Application) Export(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {

		cfg, err := app.Runit.Export(service.Name)
		if err != nil {
			log.Warnf("export %s: %s", service.Name, err)
			continue
		}

		destfile := filepath.Join("./", service.Name+".yaml")
		err = cfg.SaveFile(destfile)
		if err != nil {
			log.Warnf("export %s: %s", service.Name, err)
			continue
		}
	}
}
