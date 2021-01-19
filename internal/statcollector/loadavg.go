package statcollector

import (
	"fmt"
	"io/ioutil"
	"runtime"
)

// LoadAVG is process load average.
type LoadAVG struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// ParseLoadAVG return load average info.
func ParseLoadAVG() (*LoadAVG, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		return parseLoadAVGLinux()
	default:
		return nil, fmt.Errorf("%s OS not supported", os)
	}
}

func parseLoadAVGLinux() (*LoadAVG, error) {
	procF := ProcPath + "loadavg"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	la := LoadAVG{}
	fmt.Sscanf(string(b), "%f %f %f", &la.Load1, &la.Load5, &la.Load15)

	return &la, nil
}
