package runit

import (
	"fmt"
	"os"
)

func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Size:      1024 * 1024 * 50,
		Num:       10,
		Timestamp: 2,
	}
}

type LoggingConfig struct {
	Directory string

	// max file size in bytes
	Size int

	// number of old log files to keep. 0 means do not remove
	Num int

	// minimum number of log files
	Min int

	// empty after a period of time
	Timeout int

	Processor []string

	UdpAddress []string
	TcpAddress []string

	Prefix string

	// TODO: pattern selection

	// controls the timestamp format, 0 is disabled, 1 is -t, 2 is -tt, 3 is -ttt (see man svlogd)
	Timestamp int
}

func (cfg *LoggingConfig) WriteRunFile(path string) error {
	tmpname := path + ".tmp"
	f, err := os.Create(tmpname)
	if err != nil {
		return err
	}
	defer f.Close()

	svlogd_flags := ""

	switch cfg.Timestamp {
	case 0:
		break
	case 1:
		svlogd_flags += "-t"
	case 2:
		svlogd_flags += "-tt"
	case 3:
		svlogd_flags += "-ttt"
	default:
		log.Warnf("invalid timestamp value: %d", cfg.Timestamp)
	}

	fmt.Fprintf(f, "#!/bin/bash\n")
	fmt.Fprintf(f, "exec svlogd %s %s\n", svlogd_flags, cfg.Directory)

	err = os.MkdirAll(cfg.Directory, 0755)
	if err != nil {
		return err
	}
	f.Chmod(0755)

	err = os.Rename(tmpname, path)

	return err

}

func (cfg *LoggingConfig) WriteConfigFile(path string) error {
	tmpname := path + ".tmp"
	f, err := os.Create(tmpname)
	if err != nil {
		return err
	}
	defer f.Close()

	if cfg.Size != 0 {
		fmt.Fprintf(f, "s%d\n", cfg.Size)
	}
	if cfg.Num != 0 {
		fmt.Fprintf(f, "n%d\n", cfg.Num)
	}
	if cfg.Min != 0 {
		fmt.Fprintf(f, "N%d\n", cfg.Min)
	}
	if cfg.Timeout != 0 {
		fmt.Fprintf(f, "t%d\n", cfg.Timeout)
	}
	for _, processor := range cfg.Processor {
		fmt.Fprintf(f, "!%s\n", processor)
	}
	for _, addr := range cfg.UdpAddress {
		fmt.Fprintf(f, "u%s\n", addr)
	}
	for _, addr := range cfg.TcpAddress {
		fmt.Fprintf(f, "U%s\n", addr)
	}
	if cfg.Prefix != "" {
		fmt.Fprintf(f, "p%s\n", cfg.Prefix)
	}

	err = os.Rename(tmpname, path)

	return err
}
