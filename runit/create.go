package runit

import (
	"fmt"
)

func DefaultCreateOptions() *CreateOptions {
	return &CreateOptions{
		Script: make([]string, 0),
	}
}

type CreateOptions struct {
	Force   bool
	Restart bool
	Script  []string
}

func (opts *CreateOptions) WithScript(lines []string) *CreateOptions {
	opts.Script = append(opts.Script, lines...)
	return opts
}

func (runit *Runit) Create(cfg *ServiceConfig) error {

	return runit.Create2(cfg, nil)

}

func (runit *Runit) Create2(cfg *ServiceConfig, opts *CreateOptions) error {

	if opts == nil {
		opts = DefaultCreateOptions()
	}
	sv := runit.GetService(cfg.Name)

	if opts.Force == false && sv.Exists() {
		return fmt.Errorf("service exists")
	}

	err := runit.Apply(cfg)
	if err != nil {
		return err
	}

	if opts.Restart && sv.Running() {
		err = runit.Restart(cfg.Name)
	}

	return err

}
