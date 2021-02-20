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
		return nil, fmt.Errorf("failed of read file %s: %w", procF, err)
	}

	lc := sysmon.LoadCPU{}
	var null float64
	fmt.Sscanf(string(b), "cpu %g %g %g %g", &lc.User, &null, &lc.System, &lc.Idle)
	lc.System /= multiplier
	lc.User /= multiplier
	lc.Idle /= multiplier

	return &lc, nil
}

// AverageLoadCPU усредняет значения для массива значений LoadCPU.
func AverageLoadCPU(ll []*sysmon.LoadCPU) *sysmon.LoadCPU {
	var lu, ls, li float64
	for _, l := range ll {
		lu += l.User
		ls += l.System
		li += l.Idle
	}
	return &sysmon.LoadCPU{
		User:   lu / float64(len(ll)),
		System: ls / float64(len(ll)),
		Idle:   li / float64(len(ll)),
	}
}
