package main

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var defaultConfigString = `
# default configuration
log_directory: /var/log
`

// end default configuration

type ApplicationConfig struct {
	// where to make log files for service
	// ie: /var/log/[service]/current
	LogDirectory string `yaml:"log_directory"`

	// do sudo
	Sudo bool
}

func (c *ApplicationConfig) LoadDefault() {
	*c = *GetDefaultConfig()
}

func (c *ApplicationConfig) LoadYaml(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	if err := c.FixupConfig(); err != nil {
		return err
	}

	return nil
}

func (c *ApplicationConfig) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApplicationConfig) PrintYaml() {
	PrintConfig(c)
}

func GetDefaultConfig() *ApplicationConfig {
	cfg := &ApplicationConfig{}
	err := cfg.LoadYamlBuffer([]byte(defaultConfigString))
	if err != nil {
		panic(fmt.Sprintf("load default config: %s", err))
	}
	return cfg
}

// after loading configuration this gives us a spot to "fix up" any configuration
// or abort the loading process
func (c *ApplicationConfig) FixupConfig() error {
	// var emptyConfig ApplicationConfig

	return nil
}

func PrintDefaultConfig() {
	conf := GetDefaultConfig()
	PrintConfig(conf)
}

func PrintConfig(conf *ApplicationConfig) {
	d, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println("-- Configuration --")
	fmt.Println(string(d))
}
