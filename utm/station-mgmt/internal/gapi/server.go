package gapi

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/utm/station-mgmt/internal/service"
	smpb "github.com/utm/station-mgmt/pkg/pb/station_mgmt/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	port string
	grpc *grpc.Server
}

func NewServer(port string, overviewSvc *service.OverviewService) *Server {
	s := grpc.NewServer()
	reflection.Register(s)

	smpb.RegisterStationManagementServiceServer(s, NewStationMgmtHandler(overviewSvc))

	return &Server{port: port, grpc: s}
}

func (s *Server) Run() error {
	if s.port == "" {
		s.port = "9090"
	}
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	log.Info().Str("port", s.port).Msg("gRPC server starting")
	return s.grpc.Serve(lis)
}

func (s *Server) Stop() { s.grpc.GracefulStop() }
