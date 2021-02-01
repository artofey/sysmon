package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/artofey/sysmon/internal/pb"
	"github.com/artofey/sysmon/internal/statcollector"
)

type Server struct {
	pb.UnimplementedMonitorServer
	port   int
	StatCh chan *pb.StatSnapshot
	Stats  []*pb.StatSnapshot
}

func New(port int, st []*pb.StatSnapshot) (*Server, error) {
	return &Server{
		port:   port,
		StatCh: make(chan *pb.StatSnapshot),
		Stats:  st,
	}, nil
}

func Start(s *Server) error {
	log.Print("Start GRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMonitorServer(grpcServer, s)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *Server) GetStats(mr *pb.MonRequest, stream pb.Monitor_GetStatsServer) error {
	log.Print("new listener")

ENDFOR:
	for {
		select {
		case <-stream.Context().Done():
			log.Print("disconnect listener")
			break ENDFOR
		// Стар передачи данных после ожидания
		case <-time.After(time.Duration(mr.AveragedOver) * time.Second):
			t := time.NewTicker(time.Duration(mr.Timeout) * time.Second)
			for range t.C {
				fmt.Printf("%v - %v\n", s.Stats, mr)
				midStat, err := getMidleSnapshot(s.Stats, mr)
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

func getMidleSnapshot(st []*pb.StatSnapshot, mr *pb.MonRequest) (*pb.StatSnapshot, error) {
	var mu sync.Mutex
	ao := int(mr.AveragedOver)
	mu.Lock()
	defer mu.Unlock()

	if len(st) < ao {
		return nil, fmt.Errorf("no nedded stats")
	}

	lastSnap := st[len(st)-ao:] // get slice by last items
	var lastAVG []*pb.LoadAVG
	var lastCPU []*pb.LoadCPU

	for _, l := range lastSnap {
		lastAVG = append(lastAVG, l.Lavg)
		lastCPU = append(lastCPU, l.Lcpu)
	}

	var snap *pb.StatSnapshot
	snap.Lavg = statcollector.MidleLoadAVG(lastAVG)
	snap.Lcpu = statcollector.MidleLoadCPU(lastCPU)
	return snap, nil
}
