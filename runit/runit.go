package runit

import (
	"fmt"
	"io/ioutil"
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

func (runit *Runit) GetService(name string) *Service {
	s := &Service{
		ServiceDir: filepath.Join(runit.ServiceDir, name),
		ActiveDir:  filepath.Join(runit.ActiveDir, name),
		Name:       name,
	}
	return s
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

func (runit *Runit) Delete(name string) error {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return fmt.Errorf("no such service")
	}

	os.RemoveAll(sv.ActiveDir)
	os.RemoveAll(sv.ServiceDir)

	// TODO: cleanup logging directory?

	return nil
}

func (runit *Runit) Activate(name string) error {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return fmt.Errorf("no such service")
	}

	// activate the service
	lst, err := os.Lstat(sv.ActiveDir)
	if err != nil && os.IsNotExist(err) {
		err = os.Symlink(sv.ServiceDir, sv.ActiveDir)
	}
	if err == nil && lst.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("not a symlink: %s", sv.ActiveDir)
	}

	return nil
}

func (runit *Runit) Deactivate(name string) error {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return fmt.Errorf("no such service")
	}
	err := os.Remove(sv.ActiveDir)
	return err
}

func (runit *Runit) Disable(name string) error {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return fmt.Errorf("no such service")
	}
	if sv.Enabled() == false {
		return nil
	}
	downfile := filepath.Join(sv.ServiceDir, "down")
	err := ioutil.WriteFile(downfile, []byte{}, 0400)
	if err != nil {
		return err
	}
	return err
}
func (runit *Runit) Enable(name string) error {
	sv := runit.GetService(name)
	if sv.Exists() == false {
		return fmt.Errorf("no such service")
	}
	downfile := filepath.Join(sv.ServiceDir, "down")
	err := os.Remove(downfile)
	return err
}

func (runit *Runit) Reset(name string) error {

	return nil
}
