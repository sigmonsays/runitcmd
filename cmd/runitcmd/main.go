package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/sigmonsays/runitcmd/runit"

	gologging "github.com/sigmonsays/go-logging"
)

type Application struct {
	*cli.App
	Conf  *ApplicationConfig
	Runit *runit.Runit
}

func main() {
	c := cli.NewApp()
	c.Name = "runitcmd"
	c.Version = "0.0.1"
	c.Usage = "manage runit services"
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

		config_files := []string{
			"/etc/runitcmd.yaml",
			filepath.Join(os.Getenv("HOME"), ".runitcmd.yaml"),
		}
		app.Conf = GetDefaultConfig()
		for _, config_file := range config_files {
			_, err := os.Stat(config_file)
			if err != nil && os.IsNotExist(err) {
				continue
			}

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
		initCreate(app),
		initApply(app),
		initStatus(app),

		// more commands
		makeCommand("delete", app.Delete),
		makeCommand("activate", app.Activate),
		makeCommand("deactivate", app.Deactivate),
		makeCommand("enable", app.Enable),
		makeCommand("disable", app.Disable),
		makeCommand("reset", app.Reset),

		// commands
		makeCommand("up", app.Up),
		makeCommand("down", app.Down),
		makeCommand("pause", app.Pause),
		makeCommand("cont", app.Cont),
		makeCommand("hup", app.Cont),
		makeCommand("alarm", app.Cont),
		makeCommand("interrupt", app.Cont),
		makeCommand("quit", app.Quit),
		makeCommand("usr1", app.Usr1),
		makeCommand("usr2", app.Usr2),
		makeCommand("term", app.Term),
		makeCommand("kill", app.Kill),

		// lsb
		makeCommand("start", app.Start),
		makeCommand("stop", app.Stop),
		makeCommand("reload", app.Reload),
		makeCommand("restart", app.Restart),
		makeCommand("shutdown", app.Shutdown),
	}

	app.Run(os.Args)
}
