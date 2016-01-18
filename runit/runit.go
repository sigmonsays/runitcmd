package runit

import (
	"os"
	"path/filepath"
)

func DefaultRunitConfig() *RunitConfig {
	return &RunitConfig{
		ServiceDir: "/etc/sv",
		ActiveDir:  "/etc/service",
	}
}

type RunitConfig struct {
	ServiceDir string
	ActiveDir  string
}

func NewRunit(cfg *RunitConfig) *Runit {
	if cfg == nil {
		cfg = DefaultRunitConfig()
	}
	return &Runit{
		ServiceDir: cfg.ServiceDir,
		ActiveDir:  cfg.ActiveDir,
	}
}

type Runit struct {
	ServiceDir string
	ActiveDir  string
}

func (runit *Runit) ListServices() ([]*Service, error) {

	f, err := os.Open(runit.ServiceDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	services := make([]*Service, 0)
	for _, name := range names {
		s := &Service{
			ServiceDir: filepath.Join(runit.ServiceDir, name),
			ActiveDir:  filepath.Join(runit.ActiveDir, name),
			Name:       name,
		}
		services = append(services, s)
	}

	return services, nil
}
