package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/sigmonsays/runitcmd/runit"
)

func makeCommand(name string, fn func(*cli.Context)) cli.Command {
	cmd := cli.Command{
		Name:        name,
		Usage:       name + " a service",
		Description: name + " a service",
		Action:      fn,
	}
	return cmd
}

func (app *Application) MatchingServices(c *cli.Context) []*runit.Service {
	services := make([]*runit.Service, 0)
	args := c.Args()

	all_services, err := app.Runit.ListServices()
	if err != nil {
		log.Warnf("ListServices: %s", err)
		return nil
	}
	// TODO: Make this an option?
	use_regex := false

	var matchFn func(name string) bool
	for n := 0; n < len(args); n++ {
		pattern := args.Get(n)

		if use_regex {
			match, err := regexp.Compile(pattern)
			if err != nil {
				log.Warnf("pattern %s: %s", pattern, err)
				return nil
			}
			matchFn = func(name string) bool {
				return match.MatchString(name)
			}
		} else {
			matchFn = func(name string) bool {
				matched, _ := filepath.Match(pattern, name)
				return matched
			}
		}

		seen := make(map[string]bool, 0)
		var found bool

		for _, service := range all_services {
			if matchFn(service.Name) == false {
				continue
			}
			if _, found = seen[service.Name]; found {
				continue
			}
			services = append(services, service)
			seen[service.Name] = true
		}
	}
	return services
}

func (app *Application) runCommand(name, action string) {
	var err error

	switch action {
	case "up":
		err = app.Runit.Up(name)
	case "down":
		err = app.Runit.Down(name)
	case "once":
		err = app.Runit.Once(name)
	case "pause":
		err = app.Runit.Pause(name)
	case "cont":
		err = app.Runit.Cont(name)
	case "hup":
		err = app.Runit.Hup(name)
	case "alarm":
		err = app.Runit.Alarm(name)
	case "interrupt":
		err = app.Runit.Interrupt(name)
	case "quit":
		err = app.Runit.Quit(name)
	case "usr1":
		err = app.Runit.Usr1(name)
	case "usr2":
		err = app.Runit.Usr2(name)
	case "term":
		err = app.Runit.Term(name)
	case "kill":
		err = app.Runit.Kill(name)
	// lsb
	case "start":
		err = app.Runit.Start(name)
	case "stop":
		err = app.Runit.Stop(name)
	case "reload":
		err = app.Runit.Reload(name)
	case "restart":
		err = app.Runit.Restart(name)
	case "shutdown":
		err = app.Runit.Shutdown(name)

	default:
		err = fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		log.Warnf("%s %s: %s", name, action, err)
	}
}

func (app *Application) Up(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "up")
	}
}
func (app *Application) Down(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "down")
	}
}
func (app *Application) Once(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "once")
	}
}
func (app *Application) Pause(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "pause")
	}
}
func (app *Application) Cont(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "cont")
	}
}
func (app *Application) Hup(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "hup")
	}
}
func (app *Application) Alarm(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "alarm")
	}
}
func (app *Application) Interrupt(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "interrupt")
	}
}
func (app *Application) Quit(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "quit")
	}
}
func (app *Application) Usr1(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "1")
	}
}
func (app *Application) Usr2(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "2")
	}
}
func (app *Application) Term(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "term")
	}
}
func (app *Application) Kill(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "kill")
	}
}

// lsb compatible
func (app *Application) Start(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "start")
	}
}
func (app *Application) Stop(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "stop")
	}
}
func (app *Application) Reload(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "reload")
	}
}
func (app *Application) Restart(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "restart")
	}
}
func (app *Application) Shutdown(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "shutdown")
	}
}
