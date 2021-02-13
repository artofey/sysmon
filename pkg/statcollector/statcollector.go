// Сбор статистики и складывание ее в список

package statcollector

import (
	"fmt"
	"time"

	"github.com/artofey/sysmon"
)

type StatCollector struct {
	statsStorage []sysmon.Stats
}

func NewStatCollector() *StatCollector {
	var ss []sysmon.Stats

	c := make(chan *sysmon.Stats)
	go startColecting(c)

	for sc := range c {
		a.Stats = append(a.Stats, sc)
	
	return &StatCollector{
		statsStorage: ss,
	}
}

func (sc *StatCollector) GetStatsStorage() ([]sysmon.Stats){
	return sc.statsStorage
}

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

func getStat()(*sysmon.Stats, error) {
	avg, err := ParseLoadAVG()
	if err != nil {
		return nil, fmt.Errorf("failed parse load average: %w", err)
	}
	cpu, err := ParseLoadCPU()
	if err != nil {
		return nil, fmt.Errorf("failed parse load process: %w", err)
	}
	s := sysmon.Stats{
		Lavg: &avg,
		Lcpu: &cpu,
	}
	return &s, nil
}

func startColecting(out chan *sysmon.Stats) {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		s, err := getStat()
		if err != nil {
			break
		}
		out <- s
	}
}
