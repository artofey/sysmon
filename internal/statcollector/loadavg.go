package statcollector

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/artofey/sysmon/internal/pb"
)

// ParseLoadAVG return load average info.
func ParseLoadAVG() (*pb.LoadAVG, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		return parseLoadAVGLinux()
	default:
		return nil, fmt.Errorf("%s OS not supported", os)
	}
}

func parseLoadAVGLinux() (*pb.LoadAVG, error) {
	procF := ProcPath + "loadavg"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	la := pb.LoadAVG{}
	fmt.Sscanf(string(b), "%f %f %f", &la.Load1, &la.Load5, &la.Load15)

	return &la, nil
}
