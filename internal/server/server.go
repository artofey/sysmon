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
}

func New(port int) (*Server, error) {
	return &Server{
		StatCh: make(chan *pb.StatSnapshot),
		port:   port,
	}, nil
}

func Start(s *Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.port))
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

			if err := stream.Send(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func getStatsChan(mr *pb.MonRequest) chan *pb.StatSnapshot {
	statsC := make(chan *pb.StatSnapshot)

	return statsC
}

func getAVGSnapshot(st []*pb.StatSnapshot, mr *pb.MonRequest) (*pb.StatSnapshot, error) {
	var mu sync.Mutex
	ao := int(mr.AveragedOver)
	mu.Lock()
	defer mu.Unlock()

	if len(st) < ao {
		return nil, fmt.Errorf("no nedded stats")
	}

	lastSnap := st[len(st)-ao : len(st)] // get slice by last items
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
