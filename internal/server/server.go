package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/artofey/sysmon/internal/pb"
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

func (s *Server) GetStat(mr *pb.MonRequest, stream pb.Monitor_GetStatServer) error {
	for s := range s.StatCh {
		if err := stream.Send(s); err != nil {
			return err
		}
	}
	return nil
}
