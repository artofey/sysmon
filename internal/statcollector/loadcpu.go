package statcollector

import (
	"fmt"
	"io/ioutil"
	"runtime"
)

// LoadCPU is CPU stat.
type LoadCPU struct {
	User   uint64
	System uint64
	Idle   uint64
}

// ParseLoadCPU return CPU stat.
func ParseLoadCPU() (*LoadCPU, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		return parseLoadCPULinux()
	default:
		return nil, fmt.Errorf("%s OS not supported", os)
	}
}

func parseLoadCPULinux() (*LoadCPU, error) {
	procF := ProcPath + "stat"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	lc := LoadCPU{}
	var null uint64
	fmt.Sscanf(string(b), "cpu %d %d %d %d", &lc.User, &null, &lc.System, &lc.Idle)

	return &lc, nil
}
