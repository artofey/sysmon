package server

import (
	"log"
	"time"

	"github.com/artofey/sysmon"

	"github.com/artofey/sysmon/pkg/server/pb"
)

func (s *Server) GetStats(mr *pb.MonRequest, stream pb.Monitor_GetStatsServer) error {
	log.Print("new listener")
	defer log.Print("listener disconnected")

	for {
		select {
		case <-stream.Context().Done():
			log.Print("disconnect listener")
			return nil
		// Стар передачи данных после ожидания
		case <-time.After(time.Duration(mr.AveragedOver) * time.Second):
			t := time.NewTicker(time.Duration(mr.Timeout) * time.Second)
			for {
				select {
				case <-stream.Context().Done():
					return nil
				case <-t.C:
					midStat, err := s.sc.GetAVGStats(mrToConsumer(mr))
					if err != nil {
						return err
					}
					if err := stream.Send(statsToSnapshot(midStat)); err != nil {
						return err
					}
				}
			}
		}
	}
}

func mrToConsumer(mr *pb.MonRequest) sysmon.Consumer {
	return sysmon.Consumer{
		ID:           1,
		Timeout:      mr.GetTimeout(),
		AveragedOver: mr.GetAveragedOver(),
	}
}

func statsToSnapshot(c sysmon.Stats) *pb.StatSnapshot {
	var lavg pb.LoadAVG
	if c.Lavg == nil {
		lavg = pb.LoadAVG{}
	} else {
		lavg = pb.LoadAVG{
			Load1:  c.Lavg.Load1,
			Load5:  c.Lavg.Load5,
			Load15: c.Lavg.Load15,
		}
	}
	var lcpu pb.LoadCPU
	if c.Lcpu == nil {
		lcpu = pb.LoadCPU{}
	} else {
		lcpu = pb.LoadCPU{
			User:   c.Lcpu.User,
			System: c.Lcpu.System,
			Idle:   c.Lcpu.Idle,
		}
	}
	return &pb.StatSnapshot{
		Lavg: &lavg,
		Lcpu: &lcpu,
	}
}
