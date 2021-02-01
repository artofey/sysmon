package statcollector

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/artofey/sysmon/internal/pb"
)

// ParseLoadCPU return CPU stat.
func ParseLoadCPU() (*pb.LoadCPU, error) {
	os := runtime.GOOS
	switch os {
	case "linux":
		return parseLoadCPULinux()
	default:
		return nil, fmt.Errorf("%s OS not supported", os)
	}
}

func parseLoadCPULinux() (*pb.LoadCPU, error) {
	procF := ProcPath + "stat"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	lc := pb.LoadCPU{}
	var null uint64
	fmt.Sscanf(string(b), "cpu %d %d %d %d", &lc.User, &null, &lc.System, &lc.Idle)

	return &lc, nil
}

// MidleLoadCPU усредняет значения для массива значений LoadCPU.
func MidleLoadCPU(ll []*pb.LoadCPU) *pb.LoadCPU {
	var lu, ls, li uint64
	for _, l := range ll {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	return &pb.LoadCPU{
		User:   lu / uint64(len(ll)),
		System: ls / uint64(len(ll)),
		Idle:   li / uint64(len(ll)),
	}
}
