package detectors

import (
	"github.com/manus/upm/pkg/manager"
	"os/exec"
)

type Detector struct {
	managers []manager.Manager
}

func NewDetector() *Detector {
	return &Detector{
		managers: []manager.Manager{},
	}
}

func (d *Detector) Register(m manager.Manager) {
	d.managers = append(d.managers, m)
}

func (d *Detector) Detect() []manager.Manager {
	var available []manager.Manager
	for _, m := range d.managers {
		if m.IsAvailable() {
			available = append(available, m)
		}
	}
	return available
}

func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
