package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Host string
}

type Server struct {
	cfg      Config
	listener net.Listener
	server   *grpc.Server
}

func New(ctx context.Context, cfg Config) *Server {
	if cfg.Host == "" {
		log.Fatalf("error with gRPC starting: empty host - bad host")
	}

	listener, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		log.Fatalf("gRPC: failed to listen, error: %v", err)
	}

	return &Server{
		cfg:      cfg,
		listener: listener,
		server:   grpc.NewServer(),
	}
}

func (s *Server) Run(ctx context.Context) {
	log.Printf("gRPC: server started on %s", s.cfg.Host)
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("gRPC: failed to serve, error: %s", err)
		return
	}
}

func (s *Server) Stop() {
	s.server.Stop()
}

func (s *Server) Server() *grpc.Server {
	return s.server
}
