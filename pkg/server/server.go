package server

import (
	"fmt"
	"log"
	"net"

	"github.com/artofey/sysmon/pkg/statcollector"

	"google.golang.org/grpc"

	"github.com/artofey/sysmon/pkg/server/pb"
)

type Server struct {
	pb.UnimplementedMonitorServer
	grpcServer *grpc.Server
	sc         *statcollector.StatCollector
}

func NewServer(sc *statcollector.StatCollector) *Server {
	return &Server{
		sc: sc,
	}
}

func (s *Server) Run(port int) error {
	log.Print("Start GRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterMonitorServer(s.grpcServer, s)
	err = s.grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}
