package runit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (runit *Runit) Export(name string) (*ServiceConfig, error) {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return nil, fmt.Errorf("no such service")
	}

	configfile := filepath.Join(sv.ServiceDir, "service.yaml")

	cfg := &ServiceConfig{}

	err := cfg.LoadFile(configfile)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

// do the best we can loading service config from disk (for when we do not have a service.yaml file)
// for things we can't parse we will fill in sane defaults
func (runit *Runit) LoadFromDisk(name string) (*ServiceConfig, error) {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return nil, fmt.Errorf("not found")
	}
	cfg := &ServiceConfig{
		Name:      name,
		Activated: true,
	}

	// load the service/run file
	runfile := filepath.Join(sv.ServiceDir, "run")
	if f, err := os.Open(runfile); err == nil {
		fin := bufio.NewReader(f)
		for {
			line, err := fin.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Warnf("%s", err)
				break
			}
			sline := string(line)
			if strings.HasPrefix(sline, "exec ") {
				cfg.Exec = sline[5:]
				log.Infof("found exec %s", cfg.Exec)
			}
		}
		f.Close()

	}

	return cfg, nil
}
