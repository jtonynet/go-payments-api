package gormRepos

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"gorm.io/gorm"
)

type Merchant struct {
	gormConn database.Conn
	db       *gorm.DB
}

func NewMerchant(conn database.Conn) (port.MerchantRepository, error) {
	db, err := conn.GetDB(context.Background())
	if err != nil {
		return nil, fmt.Errorf("merchant repository failure on conn.GetDB()")
	}

	dbGorm, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("merchant repository failure to cast conn.GetDB() as gorm.DB")
	}

	return &Merchant{
		gormConn: conn,
		db:       dbGorm,
	}, nil
}

func (m *Merchant) FindByName(_ context.Context, name string) (*port.MerchantEntity, error) {
	merchantModel := gormModel.Merchant{}

	result := m.db.Preload("MCC").Where(&gormModel.Merchant{Name: name}).First(&merchantModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &port.MerchantEntity{
		Name: merchantModel.Name,
		MCC:  merchantModel.MCC.MCC,
	}, nil
}
