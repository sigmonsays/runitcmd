package main

import (
	"os/exec"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/sigmonsays/runitcmd/runit"
)

type simpleCommand struct {
	service string
	command string
}

func (s *simpleCommand) run() error {

	cmdline := []string{
		"sv",
		s.command, s.service,
	}
	log.Tracef("cmdline %s", cmdline)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()

	return err
}

func runCommand(service, command string) {
	cmd := &simpleCommand{
		service: service,
		command: command,
	}
	err := cmd.run()
	if err != nil {
		log.Warnf("%s %s: %s", command, service, err)
	}
}

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
	for n := 0; n < len(args); n++ {
		pattern := args.Get(n)

		match, err := regexp.Compile(pattern)
		if err != nil {
			log.Warnf("pattern %s: %s", pattern, err)
			return nil
		}

		for _, service := range all_services {
			if match.MatchString(service.Name) == false {
				continue
			}
			services = append(services, service)
		}
	}
	return services
}

func (app *Application) Up(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "up")
	}
}
func (app *Application) Down(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "down")
	}
}
func (app *Application) Once(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "once")
	}
}
func (app *Application) Pause(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "pause")
	}
}
func (app *Application) Cont(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "cont")
	}
}
func (app *Application) Hup(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "hup")
	}
}
func (app *Application) Alarm(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "alarm")
	}
}
func (app *Application) Interrupt(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "interrupt")
	}
}
func (app *Application) Quit(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "quit")
	}
}
func (app *Application) Usr1(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "1")
	}
}
func (app *Application) Usr2(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "2")
	}
}
func (app *Application) Term(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "term")
	}
}
func (app *Application) Kill(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "kill")
	}
}

// lsb compatible
func (app *Application) Start(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "start")
	}
}
func (app *Application) Stop(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "stop")
	}
}
func (app *Application) Reload(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "reload")
	}
}
func (app *Application) Restart(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "restart")
	}
}
func (app *Application) Shutdown(c *cli.Context) {
	for _, service := range app.MatchingServices(c) {
		runCommand(service.Name, "shutdown")
	}
}
