package redisRepos

import (
	"context"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"github.com/tidwall/gjson"
)

type Merchant struct {
	cacheConn database.InMemory

	merchantRepository port.MerchantRepository
}

func NewRedisMerchant(cacheConn database.InMemory, mRepository port.MerchantRepository) (port.MerchantRepository, error) {
	return &Merchant{
		cacheConn:          cacheConn,
		merchantRepository: mRepository,
	}, nil
}

func (m *Merchant) FindByName(_ context.Context, name string) (*port.MerchantEntity, error) {
	var mEntity *port.MerchantEntity

	merchantCached, err := m.cacheConn.Get(context.Background(), name)
	if err != nil {
		mEntity, err = m.merchantRepository.FindByName(context.Background(), name)
		if err != nil {
			return mEntity, err
		}

		defaultExpiration, err := m.cacheConn.GetDefaultExpiration(context.Background())
		if err != nil {
			return mEntity, err
		}

		err = m.cacheConn.Set(context.Background(), name, mEntity, defaultExpiration)
		if err != nil {
			return mEntity, err
		}

	} else {
		mEntity = &port.MerchantEntity{
			Name: gjson.Get(merchantCached, "Name").String(),
			MCC:  gjson.Get(merchantCached, "MCC").String(),
		}

	}

	return mEntity, nil
}
