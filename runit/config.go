package runit

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Name           string
	Exec           string
	Disabled       bool              `yaml:"disabled"`
	Activated      bool              `yaml:"activated"`
	Logging        *LoggingConfig    `yaml:"logging,omitempty"`
	RedirectStderr bool              `yaml:"redirect_stderr"`
	Env            map[string]string `yaml:"env,omitempty"`
	Export         map[string]string `yaml:"export,omitempty"`
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
