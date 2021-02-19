package statcollector

import (
	"fmt"
	"io/ioutil"

	"github.com/artofey/sysmon"
)

// ParseLoadCPU return CPU stat.
func ParseLoadCPU() (*sysmon.LoadCPU, error) {
	procF := ProcPath + "stat"
	b, err := ioutil.ReadFile(procF)
	if err != nil {
		return nil, fmt.Errorf("failed of read file %s: %v", procF, err)
	}

	lc := sysmon.LoadCPU{}
	var null uint64
	fmt.Sscanf(string(b), "cpu %d %d %d %d", &lc.User, &null, &lc.System, &lc.Idle)
	lc.System = lc.System / 1000
	lc.User = lc.User / 1000
	lc.Idle = lc.Idle / 1000

	return &lc, nil
}

// AverageLoadCPU усредняет значения для массива значений LoadCPU.
func AverageLoadCPU(ll []*sysmon.LoadCPU) *sysmon.LoadCPU {
	var lu, ls, li uint64
	for _, l := range ll {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	return &sysmon.LoadCPU{
		User:   lu / uint64(len(ll)),
		System: ls / uint64(len(ll)),
		Idle:   li / uint64(len(ll)),
	}
}
