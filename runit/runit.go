package runit

import (
	"fmt"
	"os"
	"path/filepath"
)

func DefaultRunitConfig() *RunitConfig {
	return &RunitConfig{
		ServiceDir: "/etc/sv",
		ActiveDir:  "/etc/service",
	}
}

type RunitConfig struct {
	ServiceDir string
	ActiveDir  string
}

func NewRunit(cfg *RunitConfig) *Runit {
	if cfg == nil {
		cfg = DefaultRunitConfig()
	}
	return &Runit{
		ServiceDir: cfg.ServiceDir,
		ActiveDir:  cfg.ActiveDir,
	}
}

type Runit struct {
	ServiceDir string
	ActiveDir  string
}

func (runit *Runit) GetService(name string) *Service {
	s := &Service{
		ServiceDir: filepath.Join(runit.ServiceDir, name),
		ActiveDir:  filepath.Join(runit.ActiveDir, name),
		Name:       name,
	}
	return s
}

func (runit *Runit) ListServices() ([]*Service, error) {

	f, err := os.Open(runit.ServiceDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	services := make([]*Service, 0)
	for _, name := range names {
		s := &Service{
			ServiceDir: filepath.Join(runit.ServiceDir, name),
			ActiveDir:  filepath.Join(runit.ActiveDir, name),
			Name:       name,
		}
		services = append(services, s)
	}

	return services, nil
}

type ServiceConfig struct {
	Name     string
	Exec     string
	Disabled bool
	Logging  *LoggingConfig

	// helpers
	RedirectStderr bool
	Env            map[string]string
	ExportEnv      map[string]string
}

func (runit *Runit) Create(cfg *ServiceConfig) error {
	sv := runit.GetService(cfg.Name)
	if sv.Exists() {
		return fmt.Errorf("service exists")
	}

	// create directories
	err := os.MkdirAll(sv.ServiceDir+"/log", 0755)
	if err != nil {
		return err
	}

	// write the run file
	runfile := filepath.Join(sv.ServiceDir, "run")
	f, err := os.Create(runfile)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "#!/bin/bash\n")
	if cfg.RedirectStderr {
		fmt.Fprintf(f, "exec 2>&1\n")
	}
	for k, v := range cfg.Env {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	for k, v := range cfg.ExportEnv {
		fmt.Fprintf(f, "export %s=%s\n", k, v)
	}
	fmt.Fprintf(f, "exec %s\n", cfg.Exec)
	f.Chmod(0755)
	f.Close()

	// write the log/run file
	logrun := filepath.Join(sv.ServiceDir, "log/run")

	err = cfg.Logging.WriteRunFile(logrun)
	if err != nil {
		return err
	}

	// write the log/config file
	logconfig := filepath.Join(sv.ServiceDir, "log/config")
	err = cfg.Logging.WriteConfigFile(logconfig)
	if err != nil {
		return err
	}

	if cfg.Disabled {
		return nil
	}

	// activate the service
	err = os.Symlink(sv.ServiceDir, sv.ActiveDir)

	return err

}
