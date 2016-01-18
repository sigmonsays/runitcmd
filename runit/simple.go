package runit

import (
	"os/exec"
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

func runCommand(service, command string) error {
	cmd := &simpleCommand{
		service: service,
		command: command,
	}
	err := cmd.run()
	if err != nil {
		log.Warnf("%s %s: %s", command, service, err)
		return err
	}
	return nil
}

func (runit *Runit) Up(name string) error {
	return runCommand(name, "up")
}
func (runit *Runit) Down(name string) error {
	return runCommand(name, "down")
}
func (runit *Runit) Once(name string) error {
	return runCommand(name, "once")
}
func (runit *Runit) Pause(name string) error {
	return runCommand(name, "pause")
}
func (runit *Runit) Cont(name string) error {
	return runCommand(name, "cont")
}
func (runit *Runit) Hup(name string) error {
	return runCommand(name, "hup")
}
func (runit *Runit) Alarm(name string) error {
	return runCommand(name, "alarm")
}
func (runit *Runit) Interrupt(name string) error {
	return runCommand(name, "interrupt")
}
func (runit *Runit) Quit(name string) error {
	return runCommand(name, "quit")
}
func (runit *Runit) Usr1(name string) error {
	return runCommand(name, "1")
}
func (runit *Runit) Usr2(name string) error {
	return runCommand(name, "2")
}
func (runit *Runit) Term(name string) error {
	return runCommand(name, "term")
}
func (runit *Runit) Kill(name string) error {
	return runCommand(name, "kill")
}

// lsb compatible
func (runit *Runit) Start(name string) error {
	return runCommand(name, "start")
}
func (runit *Runit) Stop(name string) error {
	return runCommand(name, "stop")
}
func (runit *Runit) Reload(name string) error {
	return runCommand(name, "reload")
}
func (runit *Runit) Restart(name string) error {
	return runCommand(name, "restart")
}
func (runit *Runit) Shutdown(name string) error {
	return runCommand(name, "shutdown")
}
