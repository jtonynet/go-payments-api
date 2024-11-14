package database

import (
	"context"
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

		/*
			TODO:
			For simplicity, I am using GORM's AutoMigrate. If time permits,
			I should migrate this solution to use the golang-migrate library,
			as it is more robust and scalable.
			See more at: https://github.com/golang-migrate/migrate
		*/
		db.AutoMigrate(&gormModel.Account{})
		db.AutoMigrate(&gormModel.Category{})
		db.AutoMigrate(&gormModel.MCC{})
		db.AutoMigrate(&gormModel.Balance{})
		db.AutoMigrate(&gormModel.Transaction{})
		db.AutoMigrate(&gormModel.Merchant{})

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
