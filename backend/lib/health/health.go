package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
)

type Config struct {
	Host          string
	Port          string
	MaxGoRoutines int
	LoggerFields  logger.Fields
}

type HealthCheckService interface {
	AddServiceReadinessHealthCheck()
	SetServiceReady(bool)
	AddPostgresReadinessCheck()
	AddPostgresReadinessCheckConnection(*pgxpool.Pool)
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type healthCheckServiceProperties struct {
	conf    Config
	logger  *logger.StandardLogger
	handler healthcheck.Handler
	server  *http.Server
}

func New(ctx context.Context, config Config, logger *logger.StandardLogger) HealthCheckService {
	s := &healthCheckServiceProperties{
		conf:    config,
		logger:  logger,
		handler: healthcheck.NewHandler(),
	}
	s.AddServiceReadinessHealthCheck()
	if config.MaxGoRoutines > 0 {
		s.addLivenessHealthCheck()
	}
	s.Start(ctx)
	return s
}

// postgresPingCheck returns a Check that validates connectivity through a postgres pool connection to a database
func (s *healthCheckServiceProperties) postgresPingCheck(conn *pgxpool.Pool, timeout time.Duration) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// lets make sure we can query database with a non-specific query
		row := conn.QueryRow(ctx, "SELECT 1")
		var oneThing int
		return row.Scan(&oneThing)
	}
}
func (s *healthCheckServiceProperties) AddPostgresReadinessCheckConnection(conn *pgxpool.Pool) {
	s.handler.AddReadinessCheck("database", healthcheck.Async(s.postgresPingCheck(conn, 1*time.Second), time.Second))
}

func (s *healthCheckServiceProperties) AddPostgresReadinessCheck() {
	s.handler.AddReadinessCheck("database", healthcheck.Async(func() error { return errors.New("no database connection") }, time.Second))
}
func (s *healthCheckServiceProperties) addLivenessHealthCheck() {
	s.handler.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(s.conf.MaxGoRoutines))
}

func (s *healthCheckServiceProperties) AddServiceReadinessHealthCheck() {
	s.handler.AddReadinessCheck("service", healthcheck.Async(func() error { return errors.New("service not started") }, time.Second))
}

func (s *healthCheckServiceProperties) SetServiceReady(ready bool) {
	s.handler.AddReadinessCheck("service", healthcheck.Async(func() error {
		if ready {
			return nil
		} else {
			return errors.New("service not ready")
		}
	}, time.Second))
}

func (s *healthCheckServiceProperties) Start(ctx context.Context) {
	address := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	s.server = &http.Server{Addr: address, Handler: s.handler}
	s.logFields(ctx).Infof("Starting health check service on %s", address)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				s.logFields(ctx).Errorf("Failed to start health check service on %s : %+v", address, err)
			}
		}
	}()
}

func (s *healthCheckServiceProperties) Stop(ctx context.Context) {
	address := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	s.logFields(ctx).Infof("Stopping health check service on %s", address)

	if err := s.server.Shutdown(ctx); err != nil {
		s.logFields(ctx).Errorf("Failed to stop health check service on %s : %+v", address, err)
	}
}

func (s *healthCheckServiceProperties) logFields(ctx context.Context) *logrus.Entry {
	return s.logger.Fields(ctx, s.conf.LoggerFields)
}
