package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/sigmonsays/runitcmd/runit"
	"github.com/urfave/cli/v2"

	gologging "github.com/sigmonsays/go-logging"
)

type Application struct {
	*cli.App
	Conf  *ApplicationConfig
	Runit *runit.Runit
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	c := cli.NewApp()
	c.Name = "runitcmd"
	c.Version = "0.0.1"
	c.Usage = "manage runit services"
	app := &Application{
		App: c,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config, c",
			Usage: "override configuration file",
		},
		&cli.StringFlag{
			Name:  "level, l",
			Value: "WARN",
			Usage: "change log level",
		},
		&cli.StringFlag{
			Name:    "log-dir",
			Usage:   "change service dir",
			EnvVars: []string{"LOG_DIR"},
			Value:   runit.DefaultLogDir,
		},
		&cli.StringFlag{
			Name:    "service-dir",
			Usage:   "change service dir",
			EnvVars: []string{"SERVICE_DIR"},
		},
		&cli.StringFlag{
			Name:    "active-dir",
			Usage:   "change active service dir",
			EnvVars: []string{"ACTIVE_SERVICE_DIR"},
		},
	}

	app.Before = func(c *cli.Context) error {
		gologging.SetLogLevel(c.String("level"))

		config_file := c.String("config")
		config_files := []string{
			"/etc/runitcmd.yaml",
		}

		user_config := filepath.Join(os.Getenv("HOME"), ".runitcmd.yaml")
		if config_file == "" {
			config_files = append(config_files, user_config)
		} else {
			config_files = []string{config_file}
		}

		app.Conf = GetDefaultConfig()
		for _, config_file := range config_files {
			_, err := os.Stat(config_file)
			if err != nil && os.IsNotExist(err) {
				continue
			}

			log.Tracef("Loading %s", config_file)
			err = app.Conf.LoadYaml(config_file)
			if err != nil {
				log.Warnf("load %s: %s", config_file, err)
			}
		}

		if app.Conf.Sudo {
			current_user, err := user.Current()
			if err == nil && current_user.Uid != "0" {
				log.Tracef("sudo uid:%s", current_user.Uid)

				argv0, err := exec.LookPath(os.Args[0])
				if err != nil {
					return err
				}

				sudo, err := exec.LookPath("sudo")
				if err != nil {
					log.Errorf("no sudo: %s", err)
					os.Exit(1)
				}

				args := []string{}
				args = append(args, sudo)
				args = append(args, "--")
				args = append(args, os.Args...)
				args[2] = argv0

				log.Tracef("sudo:%s args:%s", sudo, args)
				if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
					log.Errorf("exec: %s", err)
					os.Exit(1)
				}
				// never reached
				panic("what")
				os.Exit(0)

			}

		}

		rcfg := runit.DefaultRunitConfig()
		service_dir := c.String("service-dir")
		active_dir := c.String("active-dir")
		if service_dir == "" {
			rcfg.ServiceDir = app.Conf.ServiceDir
		} else {
			rcfg.ServiceDir = service_dir
		}
		if active_dir == "" {
			rcfg.ActiveDir = app.Conf.ActiveDir
		} else {
			rcfg.ActiveDir = active_dir
		}

		app.Runit = runit.NewRunit(rcfg)

		return nil
	}

	app.Commands = []*cli.Command{
		initList(app),
		initCreate(app),
		initSetup(app),
		initApply(app),
		initStatus(app),
		initImport(app),
		initExport(app),

		// more commands
		makeCommand(app, "delete", "", "delete service"),
		makeCommand(app, "activate", "", "create service symlink"),
		makeCommand(app, "deactivate", "", "delete service symlink"),
		makeCommand(app, "enable", "", "enable service at boot"),
		makeCommand(app, "disable", "", "disable service at boot"),
		makeCommand(app, "reset", "", "reset service state"),

		// commands
		makeCommand(app, "up", "", "bring service up"),
		makeCommand(app, "down", "", "bring service down"),
		makeCommand(app, "pause", "", "pause service"),

		// signals
		makeCommand(app, "cont", "", "send service CONT signal"),
		makeCommand(app, "hup", "", "send service HUP signal"),
		makeCommand(app, "alarm", "", "send service ALRM signal"),
		makeCommand(app, "interrupt", "", "send service INT signal"),
		makeCommand(app, "quit", "", "send service QUIT signal"),
		makeCommand(app, "usr1", "1", "send service USR1 signal"),
		makeCommand(app, "usr2", "2", "send service USR2 signal"),
		makeCommand(app, "term", "", "send service TERM signal"),
		makeCommand(app, "kill", "", "send service KILL signal"),

		// lsb
		makeCommand(app, "start", "", "start service"),
		makeCommand(app, "stop", "", "stop service"),
		makeCommand(app, "reload", "", "reload service config"),
		makeCommand(app, "restart", "", "restsart service"),
		makeCommand(app, "shutdown", "", "shutdown service"),
	}

	app.Run(os.Args)
}
