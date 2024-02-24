package main

import (
	"fmt"
	"path/filepath"

	"github.com/sigmonsays/runitcmd/runit"
	"github.com/urfave/cli/v2"

	gologging "github.com/sigmonsays/go-logging"
)

// for compatibility and legacy; use create instead
func initSetup(app *Application) *cli.Command {
	description := "setup a service"
	usage := "setup a service"

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "service name",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "be verbose",
		},
		&cli.BoolFlag{
			Name:    "enable",
			Aliases: []string{"e"},
			Usage:   "setup service and enable it",
		},
		&cli.BoolFlag{
			Name:    "disable",
			Aliases: []string{"d"},
			Usage:   "setup service but do not enable it",
		},
		&cli.StringFlag{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run command",
		},
		&cli.StringFlag{
			Name:  "log-run",
			Usage: "service log command",
		},
		&cli.StringSliceFlag{
			Name:    "script",
			Aliases: []string{"s"},
			Usage:   "additional lines to execute before run",
		},
		&cli.BoolFlag{
			Name:  "restart",
			Usage: "restart service",
		},
		&cli.IntFlag{
			Name:    "uid",
			Aliases: []string{"u"},
			Usage:   "user id",
		},
		&cli.IntFlag{
			Name:    "gid",
			Aliases: []string{"g"},
			Usage:   "group id",
		},
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "service template to use (does nothing, legacy only)",
		},
	}

	cmd := &cli.Command{
		Name:        "setup",
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      app.Setup,
	}
	return cmd
}

func (app *Application) Setup(c *cli.Context) error {
	verbose := c.Bool("verbose")
	service_dir := c.String("service-dir")
	active_dir := c.String("active-dir")
	log_dir := c.String("log-dir")
	enable := c.Bool("enable")
	disable := c.Bool("disable")
	run := c.String("run")
	log_run := c.String("log-run")

	script := c.StringSlice("script")
	restart := c.Bool("restart")
	uid := c.Int("uid")
	gid := c.Int("gid")

	args := c.Args()
	name := args.First()
	if name == "" {
		name = c.String("name")
	}

	log.Tracef("setup service-dir:%s active-dir:%s log-dir:%s",
		service_dir, active_dir, log_dir)

	if verbose {
		gologging.SetLogLevel("trace")
	}
	// template does nothing

	// make the runit api locally for these flags to function
	rcfg := runit.DefaultRunitConfig()
	if service_dir != "" {
		rcfg.ServiceDir = service_dir
	}
	if active_dir != "" {
		rcfg.ActiveDir = active_dir
	}
	rapi := runit.NewRunit(rcfg)

	if name == "" {
		log.Errorf("service name is required")
		return fmt.Errorf("service name is required")
	}
	if run == "" {
		log.Errorf("run command is required")
		return fmt.Errorf("run command is required")
	}

	if log_run != "" {
		log.Warnf("specifying log-run is not supported, argument is ignored")
	}

	log.Tracef("setup %s", name)

	lcfg := runit.DefaultLoggingConfig()
	lcfg.Directory = filepath.Join(log_dir, name)

	var exec string
	if uid != 0 || gid != 0 {
		exec = fmt.Sprintf("chpst -u %d:%d %s", uid, gid, run)
	} else {
		exec = run
	}

	if enable {
		disable = false
	}

	// log_run is not supported

	cfg := &runit.ServiceConfig{
		Name:           name,
		Exec:           exec,
		Logging:        lcfg,
		Disabled:       disable,
		Activated:      true,
		RedirectStderr: true,
	}

	copts := &runit.CreateOptions{
		Force:   true,
		Restart: restart,
		Script:  script,
	}

	err := rapi.Create2(cfg, copts)
	if err != nil {
		log.Warnf("Create2 %s: %s", name, err)
	}

	return err
}
