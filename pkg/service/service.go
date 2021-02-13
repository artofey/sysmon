package service

import (
	"github.com/artofey/sysmon"
	"github.com/artofey/sysmon/pkg/statcollector"
)

type StatCollector interface {
	GetStat() (*sysmon.Stats, error)
	StartColecting(chan *sysmon.Stats)
}

type Service struct {
	StatCollector
}

func NewService() *Service {
	return &Service{
		StatCollector: statcollector.NewStatCollector(),
	}
}
