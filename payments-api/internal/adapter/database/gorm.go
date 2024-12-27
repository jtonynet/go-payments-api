package database

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

type GormConn struct {
	db       *gorm.DB
	strategy string
	driver   string
}

func NewGormConn(cfg config.Database) (Conn, error) {
	switch cfg.Driver {
	case "postgres":
		strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host,
			cfg.User,
			cfg.Pass,
			cfg.DB,
			cfg.Port,
			cfg.SSLmode)
		db, err := gorm.Open(postgres.Open(strConn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failure on database connection: %w", err)
		}

		if cfg.MetricEnabled {
			db.Use(prometheus.New(prometheus.Config{
				DBName:          cfg.MetricDBName,        // `DBName` as metrics label
				RefreshInterval: cfg.MetricIntervalInSec, // refresh metrics interval (default 15 seconds)
				StartServer:     cfg.MetricStartServer,   // start http server to expose metrics
				HTTPServerPort:  cfg.MetricServerPort,    // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
				MetricsCollector: []prometheus.MetricsCollector{
					&prometheus.Postgres{VariableNames: []string{"Threads_running"}},
				},
			}))
		}

		gConn := GormConn{
			db:       db,
			strategy: cfg.Strategy,
			driver:   cfg.Driver,
		}

		return gConn, nil

	default:
		return nil, fmt.Errorf("database conn driver not suported: %s", cfg.Driver)
	}
}

func (gConn GormConn) Readiness(_ context.Context) error {
	rawDB, err := gConn.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := rawDB.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}

	return nil
}

func (gConn GormConn) GetDB(_ context.Context) (interface{}, error) {
	return gConn.db, nil
}

func (gConn GormConn) GetStrategy(_ context.Context) (string, error) {
	return gConn.strategy, nil
}

func (gConn GormConn) GetDriver(_ context.Context) (string, error) {
	return gConn.driver, nil
}
