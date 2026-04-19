package gapi

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/utm/backend/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg  *config.Config
	grpc *grpc.Server
}

func NewServer(cfg *config.Config) *Server {
	s := grpc.NewServer()
	reflection.Register(s) // cho phép grpcurl introspect

	// TODO: đăng ký các gRPC service tại đây
	// pb.RegisterUTMServiceServer(s, &UTMServiceHandler{})

	return &Server{cfg: cfg, grpc: s}
}

func (s *Server) Run() error {
	port := s.cfg.App.GRPCPort
	if port == "" {
		port = "9090"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	log.Info().Str("port", port).Msg("gRPC server starting")
	return s.grpc.Serve(lis)
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
