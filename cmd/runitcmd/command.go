package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/sigmonsays/runitcmd/runit"
	"github.com/urfave/cli/v2"
)

func makeCommand(name string, fn func(*cli.Context) error) *cli.Command {
	cmd := &cli.Command{
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
	for n := 0; n < args.Len(); n++ {
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

func (app *Application) runCommand(name, action string) error {
	var err error

	sv := app.Runit.GetService(name)
	if sv.Exists() == false {
		log.Warnf("GetService %s: no such service", name)
		return fmt.Errorf("GetService %s: no such service", name)

	}

	switch action {

	// commands
	case "delete":
		err = app.Runit.Delete(name)
	case "activate":
		err = app.Runit.Activate(name)
	case "deactivate":
		err = app.Runit.Deactivate(name)
	case "disable":
		err = app.Runit.Disable(name)
	case "enable":
		err = app.Runit.Enable(name)
	case "reset":
		err = app.Runit.Reset(name)

	// sv tasks
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
	return err
}

// some helpful commands
func (app *Application) Delete(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "delete")
	}
	return nil
}
func (app *Application) Activate(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "activate")
	}
	return nil
}
func (app *Application) Deactivate(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "deactivate")
	}
	return nil
}
func (app *Application) Enable(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "enable")
	}
	return nil
}
func (app *Application) Disable(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "disable")
	}
	return nil
}
func (app *Application) Reset(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "reset")
	}
	return nil
}

// sv commands
func (app *Application) Up(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "up")
	}
	return nil
}
func (app *Application) Down(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "down")
	}
	return nil
}
func (app *Application) Once(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "once")
	}
	return nil
}
func (app *Application) Pause(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "pause")
	}
	return nil
}
func (app *Application) Cont(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "cont")
	}
	return nil
}
func (app *Application) Hup(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "hup")
	}
	return nil
}
func (app *Application) Alarm(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "alarm")
	}
	return nil
}
func (app *Application) Interrupt(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "interrupt")
	}
	return nil
}
func (app *Application) Quit(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "quit")
	}
	return nil
}
func (app *Application) Usr1(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "1")
	}
	return nil
}
func (app *Application) Usr2(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "2")
	}
	return nil
}
func (app *Application) Term(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "term")
	}
	return nil
}
func (app *Application) Kill(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "kill")
	}
	return nil
}

// lsb compatible
func (app *Application) Start(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "start")
	}
	return nil
}
func (app *Application) Stop(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "stop")
	}
	return nil
}
func (app *Application) Reload(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "reload")
	}
	return nil
}
func (app *Application) Restart(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "restart")
	}
	return nil
}
func (app *Application) Shutdown(c *cli.Context) error {
	for _, service := range app.MatchingServices(c) {
		app.runCommand(service.Name, "shutdown")
	}
	return nil
}
