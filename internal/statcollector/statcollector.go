package statcollector

import (
	"fmt"
	"time"
)

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

// StatCollector сontains all types of statistics.
type StatCollector struct {
	Lavg *LoadAVG
	Lcpu *LoadCPU
}

func New() *StatCollector {
	return &StatCollector{
		Lavg: &LoadAVG{},
		Lcpu: &LoadCPU{},
	}
}

func GetStat() (*StatCollector, error) {
	avg, err := ParseLoadAVG()
	if err != nil {
		return nil, fmt.Errorf("failed parse load average: %w", err)
	}
	cpu, err := ParseLoadCPU()
	if err != nil {
		return nil, fmt.Errorf("failed parse load process: %w", err)
	}
	s := StatCollector{
		Lavg: avg,
		Lcpu: cpu,
	}
	return &s, nil
}

func StartColecting(out chan *StatCollector) {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		s, err := GetStat()
		if err != nil {
			break
		}
		out <- s
	}
}
