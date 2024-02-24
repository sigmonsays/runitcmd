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
	fmt.Fprintf(f, "#!/usr/bin/env bash\n")
	if cfg.RedirectStderr {
		fmt.Fprintf(f, "exec 2>&1\n")
	}
	for k, v := range cfg.Env {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	for k, v := range cfg.Export {
		fmt.Fprintf(f, "export %s=%s\n", k, v)
	}
	for _, script := range cfg.InlineScript {
		fmt.Fprint(f, script)
	}

	if cfg.Script != "" {
		fmt.Fprintf(f, cfg.Script)
	}

	fmt.Fprintf(f, "exec %s\n", cfg.Exec)
	f.Chmod(0755)
	f.Close()

	log.Tracef("run exec %s", cfg.Exec)

	// write the log/run file
	logrun := filepath.Join(sv.ServiceDir, "log/run")

	if cfg.Logging == nil {
		cfg.Logging = &LoggingConfig{
			Directory: filepath.Join(DefaultLogDir, sv.Name),
		}
	}

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

	// write the configuration down
	configfile := filepath.Join(sv.ServiceDir, "service.yaml")
	err = cfg.SaveFile(configfile)
	if err != nil {
		return err
	}

	if cfg.Disabled {
		err = runit.Deactivate(sv.Name)
		if err != nil {
			return err
		}
	}
	if cfg.Activated {

		err = runit.Activate(sv.Name)
		if err != nil {
			return err
		}

	}

	return err

}
