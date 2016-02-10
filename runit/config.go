package runit

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Name           string
	Exec           string
	Disabled       bool              `yaml:"disabled,omitempty"`
	Activated      bool              `yaml:"activated"`
	Logging        *LoggingConfig    `yaml:"logging,omitempty"`
	RedirectStderr bool              `yaml:"redirect_stderr,omitempty"`
	Env            map[string]string `yaml:"env,omitempty"`
	Export         map[string]string `yaml:"export,omitempty"`

	InlineScript []string `yaml:"-"`

	Script string `yaml:"script,omitempty"`
}

func (c *ServiceConfig) LoadFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ServiceConfig) SaveFile(path string) error {
	tmppath := path + ".tmp"

	data, err := yaml.Marshal(c)
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(tmppath, data, 0644)
	if err != nil {
		return nil
	}

	err = os.Rename(tmppath, path)
	return err
}

func (s *Service) Config() (*ServiceConfig, error) {
	cfgfile := filepath.Join(s.ServiceDir, "service.yaml")
	c := &ServiceConfig{}
	err := c.LoadFile(cfgfile)
	if err != nil {
		return nil, err
	}
	return c, nil
}
