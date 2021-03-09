// Сбор статистики и ее выдача

package statcollector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artofey/sysmon"
)

type StorageI interface {
	Add(sysmon.Stats) error
	GetLast(int) []sysmon.Stats
	Len() int
}

type Parser interface {
	Parse() (interface{}, error)
	Average(interface{}) (interface{}, error)
	Valid(interface{}) bool
}

type StatCollector struct {
	Parsers []Parser
	Storage StorageI
}

func NewStatCollector(s StorageI, pp []Parser) *StatCollector {
	return &StatCollector{
		Parsers: pp,
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
		stats, err := s.parseAllStats()
		if err != nil {
			log.Print(err.Error())
			return
		}
		if err = s.Storage.Add(stats); err != nil {
			log.Print(err.Error())
			return
		}
	}
}

func (s *StatCollector) AVGStats(consumer sysmon.Consumer) (sysmon.Stats, error) {
	var nilStats sysmon.Stats
	ao := int(consumer.AveragedOver)

	if s.Storage.Len() < ao {
		return nilStats, fmt.Errorf("no nedded stats")
	}

	lastStats := s.Storage.GetLast(ao)
	lastAVG := make([]*sysmon.LoadAVG, 0, ao)
	lastCPU := make([]*sysmon.LoadCPU, 0, ao)

	for _, l := range lastStats {
		lastAVG = append(lastAVG, l.Lavg)
		lastCPU = append(lastCPU, l.Lcpu)
	}

	var avgStats sysmon.Stats
	for _, p := range s.Parsers {
		switch {
		case p.Valid(lastAVG[0]):
			avg, err := p.Average(lastAVG)
			if err != nil {
				return nilStats, fmt.Errorf("average error: %w", err)
			}
			avgStats.Lavg = avg.(*sysmon.LoadAVG)
		case p.Valid(lastCPU[0]):
			avg, err := p.Average(lastCPU)
			if err != nil {
				return nilStats, fmt.Errorf("average error: %w", err)
			}
			avgStats.Lcpu = avg.(*sysmon.LoadCPU)
		}
	}
	return avgStats, nil
}

func (s *StatCollector) parseAllStats() (sysmon.Stats, error) {
	newStats := sysmon.Stats{}
	for _, parser := range s.Parsers {
		stat, err := parser.Parse()
		if err != nil {
			return sysmon.Stats{}, err
		}
		switch stat.(type) {
		case *sysmon.LoadAVG:
			newStats.Lavg = stat.(*sysmon.LoadAVG)
		case *sysmon.LoadCPU:
			newStats.Lcpu = stat.(*sysmon.LoadCPU)
		default:
			return sysmon.Stats{}, fmt.Errorf("incorrect parser type")
		}
	}
	return newStats, nil
}
