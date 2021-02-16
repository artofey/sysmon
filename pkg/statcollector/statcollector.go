// Сбор статистики и складывание ее в список

package statcollector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artofey/sysmon"
)

// ProcPath is path to proc dir.
var ProcPath string = "/proc/"

type StorageI interface {
	Add(sysmon.Stats) error
	GetLast(int) []sysmon.Stats
	Len() int
}

type StatCollector struct {
	Storage StorageI
}

func NewStatCollector(s StorageI) *StatCollector {
	return &StatCollector{
		Storage: s,
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
			log.Print(err.Error())
			return
		}
		if err = s.Storage.Add(stat); err != nil {
			log.Print(err.Error())
			return
		}
	}
}

func (s *StatCollector) GetAVGStats(consumer sysmon.Consumer) (sysmon.Stats, error) {
	var nilStats sysmon.Stats
	ao := int(consumer.AveragedOver)

	if s.Storage.Len() < ao {
		return nilStats, fmt.Errorf("no nedded stats")
	}

	lastSnap := s.Storage.GetLast(ao)
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
