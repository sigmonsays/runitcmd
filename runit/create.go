package runit

import (
	"fmt"
)

func (runit *Runit) Create(cfg *ServiceConfig) error {
	sv := runit.GetService(cfg.Name)
	if sv.Exists() {
		return fmt.Errorf("service exists")
	}

	err := runit.Apply(cfg)
	return err

}
