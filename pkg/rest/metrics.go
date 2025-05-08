package rest

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	Port   int
	server *fiber.App
}

func NewMetricsServer(port int) *MetricsServer {
	if port == 0 {
		port = 4000
	}

	s := fiber.New(fiberConfig)

	s.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return &MetricsServer{
		Port:   port,
		server: s,
	}
}

func (s *MetricsServer) Run(ctx context.Context) {
	log.Printf("metrics server was started on :%d", s.Port)
	if err := s.server.Listen(fmt.Sprintf(":%d", s.Port)); err != nil {
		log.Fatal("metrics server was stopped", err)
	}
}

func (s *MetricsServer) Shutdown() error {
	return s.server.Shutdown()
}
