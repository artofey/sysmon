package server

import (
	"log"
	"time"

	"github.com/artofey/sysmon"

	"github.com/artofey/sysmon/pkg/server/pb"
)

func (s *Server) GetStats(mr *pb.MonRequest, stream pb.Monitor_GetStatsServer) error {
	log.Print("new listener")

	for {
		select {
		case <-stream.Context().Done():
			log.Print("disconnect listener")
			return nil
		// Стар передачи данных после ожидания
		case <-time.After(time.Duration(mr.AveragedOver) * time.Second):
			t := time.NewTicker(time.Duration(mr.Timeout) * time.Second)
			for range t.C {
				// fmt.Printf("%v - %v\n", s.Stats, mr)
				midStat, err := s.sc.GetAVGStats(mrToConsumer(mr))
				if err != nil {
					return err
				}
				if err := stream.Send(midStat); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func mrToConsumer(mr *pb.MonRequest) sysmon.Consumer {
	return sysmon.Consumer{
		ID:           1,
		Timeout:      mr.GetTimeout(),
		AveragedOver: mr.GetAveragedOver(),
	}
}

func statsToSnapshot(c sysmon.Stats) *pb.StatSnapshot {
	return &pb.StatSnapshot{
		// TODO
	}
}
