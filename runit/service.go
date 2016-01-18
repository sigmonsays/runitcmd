package runit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Status int

const (
	StatusInvalid Status = iota
	StatusUp
	StatusDown
)

const svTimeMod = 4611686018427387914

func svNow() uint64 {
	return uint64(svTimeMod + time.Now().Unix())
}

type ServiceStatus struct {
	Enabled bool
	Running bool
	Pid     int
	Seconds int64
}

type Service struct {
	ActiveDir  string
	ServiceDir string
	Name       string
}

func (s *Service) String() string {
	return fmt.Sprintf("service(name=%s dir=%s/%s)", s.Name, s.ActiveDir, s.ServiceDir)
}

func (s *Service) Exists() bool {
	runfile := filepath.Join(s.ServiceDir, "run")
	_, err := os.Stat(runfile)

	if err == nil {
		return true
	}
	return false
}

func (s *Service) Enabled() bool {
	_, err := os.Stat(s.ActiveDir)
	if err == nil {
		return true
	}
	return false
}

func (s *Service) Running() bool {
	statfile := filepath.Join(s.ActiveDir, "supervise/stat")
	data, err := ioutil.ReadFile(statfile)
	if err != nil {
		return false
	}

	stat := strings.TrimRight(string(data), "\n")
	if stat == "run" {
		return true
	}
	return false

}

func (s *Service) Status() (*ServiceStatus, error) {
	st := &ServiceStatus{
		Enabled: s.Enabled(),
		Running: s.Running(),
	}

	statusfile := filepath.Join(s.ActiveDir, "supervise/status")

	if _, err := os.Stat(statusfile); err != nil && os.IsNotExist(err) {
		return st, nil
	}

	status, err := ioutil.ReadFile(statusfile)
	if err != nil {
		return nil, err
	}

	pid := uint(status[15])
	pid <<= 8
	pid += uint(status[14])
	pid <<= 8
	pid += uint(status[13])
	pid <<= 8
	pid += uint(status[12])
	st.Pid = int(pid)

	t := uint64(status[0])
	t <<= 8
	t += uint64(status[1])
	t <<= 8
	t += uint64(status[2])
	t <<= 8
	t += uint64(status[3])
	t <<= 8
	t += uint64(status[4])
	t <<= 8
	t += uint64(status[5])
	t <<= 8
	t += uint64(status[6])
	t <<= 8
	t += uint64(status[7])
	st.Seconds = int64(svNow() - t)

	return st, nil
}
