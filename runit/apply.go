package runit

import (
	"fmt"
	"os"
	"path/filepath"
)

func (runit *Runit) Apply(cfg *ServiceConfig) error {
	sv := runit.GetService(cfg.Name)

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
	for k, v := range cfg.Export {
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

	if cfg.Logging != nil {
		// write the log/config file
		logconfig := filepath.Join(sv.ServiceDir, "log/config")
		err = cfg.Logging.WriteConfigFile(logconfig)
		if err != nil {
			return err
		}
	}

	// write the configuration down
	configfile := filepath.Join(sv.ServiceDir, "service.yaml")
	err = cfg.SaveFile(configfile)
	if err != nil {
		return err
	}

	if cfg.Disabled == false {

		err = runit.Deactivate(sv.Name)
		if err != nil {
			return err
		}
	}
	if cfg.Activated == false {

		err = runit.Activate(sv.Name)
		if err != nil {
			return err
		}

	}

	return err

}
