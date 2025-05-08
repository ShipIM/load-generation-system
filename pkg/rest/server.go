package rest

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Run(ctx context.Context)
	Shutdown(ctx context.Context) error
	Router() *fiber.App
}

type Config struct {
	Host        string
	MetricsPort int
}

type server struct {
	app           *fiber.App
	metricsServer *MetricsServer
	config        Config
}

var fiberConfig = fiber.Config{
	DisableStartupMessage:    true,
	CaseSensitive:            true,
	StrictRouting:            false,
	EnableSplittingOnParsers: true,
}

func New(cfg Config) (srv Server) {
	s := &server{
		app:           fiber.New(fiberConfig),
		config:        cfg,
		metricsServer: NewMetricsServer(cfg.MetricsPort),
	}

	return s
}

func (s *server) Run(ctx context.Context) {
	go s.metricsServer.Run(ctx)

	log.Printf("server was started on %s", s.config.Host)
	if err := s.app.Listen(s.config.Host); err != nil {
		log.Fatalf("server was stopped: %s", err.Error())
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.metricsServer.Shutdown(); err != nil {
		log.Printf("failed to stop metrics server: %v", err)
	}

	return s.app.Shutdown()
}

func (s *server) Router() *fiber.App {
	return s.app
}
