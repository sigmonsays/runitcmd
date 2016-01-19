package runit

import (
	"fmt"
	"path/filepath"
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
