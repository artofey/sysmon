// Сбор статистики и складывание ее в список

package statcollector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/artofey/sysmon"
)

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

type StatCollector struct {
	Storage []sysmon.Stats
}

func NewStatCollector() *StatCollector {
	return &StatCollector{
		Storage: make([]sysmon.Stats, 0, 36000),
	}
}

func (s *StatCollector) StartCollecting(ctx context.Context) {
	t := time.NewTicker(1 * time.Second)

	for range t.C {
		select {
		case <-ctx.Done():
			return
		default:
		}
		stat, err := getStat()
		if err != nil {
			return
		}
		s.Storage = append(s.Storage, stat)
	}
}

func (s *StatCollector) GetAVGStats(consumer sysmon.Consumer) (sysmon.Stats, error) {
	var nilStats sysmon.Stats
	var mu sync.Mutex
	ao := int(consumer.AveragedOver)
	mu.Lock()
	defer mu.Unlock()

	if len(s.Storage) < ao {
		return nilStats, fmt.Errorf("no nedded stats")
	}

	lastSnap := s.Storage[len(s.Storage)-ao:] // get slice by last items
	var lastAVG []*sysmon.LoadAVG
	var lastCPU []*sysmon.LoadCPU

	for _, l := range lastSnap {
		lastAVG = append(lastAVG, l.Lavg)
		lastCPU = append(lastCPU, l.Lcpu)
	}

	var snap sysmon.Stats
	snap.Lavg = AverageLoadAVG(lastAVG)
	snap.Lcpu = AverageLoadCPU(lastCPU)
	return snap, nil
}

func getStat() (sysmon.Stats, error) {
	var nilStats sysmon.Stats
	avg, err := ParseLoadAVG()
	if err != nil {
		return nilStats, fmt.Errorf("failed parse load average: %w", err)
	}
	cpu, err := ParseLoadCPU()
	if err != nil {
		return nilStats, fmt.Errorf("failed parse load process: %w", err)
	}
	s := sysmon.Stats{
		Lavg: avg,
		Lcpu: cpu,
	}
	return s, nil
}
