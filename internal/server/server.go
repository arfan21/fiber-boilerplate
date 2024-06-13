package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arfan21/fiber-boilerplate/config"
	_ "github.com/arfan21/fiber-boilerplate/docs"
	"github.com/arfan21/fiber-boilerplate/pkg/exception"
	"github.com/arfan21/fiber-boilerplate/pkg/logger"
	"github.com/arfan21/fiber-boilerplate/pkg/middleware"
	"github.com/arfan21/fiber-boilerplate/pkg/pkgutil"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	ctxTimeout = 5
)

type Server struct {
	app     *fiber.App
	db      *pgxpool.Pool
	dbRedis *redis.Client
}

func New(
	db *pgxpool.Pool,
	dbRedis *redis.Client,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: exception.FiberErrorHandler,
	})
	timeout := time.Duration(config.Get().Service.Timeout) * time.Second
	app.Use(middleware.Timeout(timeout))

	app.Use(cors.New())
	app.Use(otelfiber.Middleware())
	app.Use(middleware.TraceID())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		// Logger: logger.Log(context.Background()),
		GetLogger: func(c *fiber.Ctx) zerolog.Logger {
			return *logger.Log(c.UserContext())
		},
	}))

	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:     app,
		db:      db,
		dbRedis: dbRedis,
	}
}

func (s *Server) Run() error {
	s.Routes()
	ctx := context.Background()
	go func() {
		if err := s.app.Listen(pkgutil.GetPort()); err != nil {
			logger.Log(ctx).Fatal().Err(err).Msg("failed to start server")
		}
	}()

	go func() {
		logger.Log(ctx).Info().Msgf("Starting prometheus exporter on port %s", config.Get().Otel.ExporterPrometheusPort)
		http.Handle(config.Get().Otel.ExporterPrometheusPath, promhttp.Handler())
		if err := http.ListenAndServe(pkgutil.GetPort(config.Get().Otel.ExporterPrometheusPort), nil); err != nil {
			logger.Log(ctx).Fatal().Err(err).Msg("failed to start prometheus exporter")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	_, shutdown := context.WithTimeout(ctx, ctxTimeout*time.Second)
	defer shutdown()

	logger.Log(ctx).Info().Msg("shutting down server")
	return s.app.Shutdown()
}
