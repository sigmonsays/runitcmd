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
			Name:  "log-level, l",
			Usage: "log level",
			Value: "warn",
		},
		&cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "be verbose",
		},
		&cli.StringFlag{
			Name:  "service-dir",
			Usage: "service directory",
			Value: runit.DefaultServiceDir,
		},
		&cli.StringFlag{
			Name:  "active-service-dir",
			Usage: "active service directory",
			Value: runit.DefaultActiveDir,
		},
		&cli.StringFlag{
			Name:  "log-dir",
			Usage: "log to directory",
		},
		&cli.BoolFlag{
			Name:  "enable, e",
			Usage: "setup service and enable it",
		},
		&cli.BoolFlag{
			Name:  "disable, d",
			Usage: "setup service but do not enable it",
		},
		&cli.StringFlag{
			Name:  "run",
			Usage: "run command",
		},
		&cli.StringFlag{
			Name:  "log-run",
			Usage: "service log command",
		},
		&cli.StringSliceFlag{
			Name:  "script",
			Usage: "additional lines to execute before run",
		},
		&cli.BoolFlag{
			Name:  "restart, r",
			Usage: "restart service",
		},
		&cli.IntFlag{
			Name:  "uid, u",
			Usage: "user id",
		},
		&cli.IntFlag{
			Name:  "gid, g",
			Usage: "group id",
		},
		&cli.StringFlag{
			Name:  "template, t",
			Usage: "service template to use (does nothing, legacy only)",
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
	log_level := c.String("log-level")
	verbose := c.Bool("verbose")
	service_dir := c.String("service-dir")
	active_dir := c.String("active-service-dir")
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

	if verbose {
		gologging.SetLogLevel("trace")
	}
	if log_level != "" {
		gologging.SetLogLevel(log_level)
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
	if log_dir == "" {
		lcfg.Directory = filepath.Join(runit.DefaultLogDir, name)
	} else {
		lcfg.Directory = log_dir
	}

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
