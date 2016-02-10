package main

import (
	"github.com/sigmonsays/runitcmd/runit"

	"github.com/codegangsta/cli"
)

func initImport(app *Application) cli.Command {
	description := "import service"
	usage := "import service"

	flags := []cli.Flag{}

	cmd := cli.Command{
		Name:        "import",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Import,
	}
	return cmd
}

func (app *Application) Import(c *cli.Context) {
	filenames := c.Args()

	for _, filename := range filenames {
		cfg := &runit.ServiceConfig{}
		err := cfg.LoadFile(filename)
		if err != nil {
			log.Warnf("import %s: %s", filename, err)
			continue
		}
		log.Debugf("loaded %s configuration", filename)
		log.Tracef("about to apply config %+v", cfg)

		err = app.Runit.Apply(cfg)
		if err != nil {
			log.Warnf("apply %s: %s", cfg.Name, err)
			continue
		}

	}
}
