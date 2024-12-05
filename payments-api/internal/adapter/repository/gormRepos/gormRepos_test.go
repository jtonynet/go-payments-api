package gormRepos

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	postgresMigrationsPath = "file://./../../database/postgres/migrations"
	postgresSeedPath       = "./../../database/postgres/seeds/integration_test_charge.up.sql"

	accountUID, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	merchantNameToMap       = "UBER EATS                   SAO PAULO BR"
	merchantCorrectMccToMap = "5412"
	merchantCategoryToMap   = uint(2)
)

type RepositoriesSuite struct {
	suite.Suite

	AccountRepo  port.AccountRepository
	MerchantRepo port.MerchantRepository

	AccountEntity port.AccountEntity
	BalanceEntity port.BalanceEntity

	migrate *migrate.Migrate
}

func (suite *RepositoriesSuite) SetupSuite() {
	cfg, err := config.LoadConfig("./../../../../")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	conn, err := database.NewConn(cfg.Database)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	if conn.Readiness(context.Background()) != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	account, err := NewGormAccount(conn)
	if err != nil {
		log.Fatalf("error when instantiating account repository: %v", err)
	}

	merchant, err := NewMerchant(conn)
	if err != nil {
		log.Fatalf("error when instantiating merchant repository: %v", err)
	}

	suite.AccountRepo = account
	suite.MerchantRepo = merchant

	suite.loadDBtestData(conn)
}

func (suite *RepositoriesSuite) loadDBtestData(conn database.Conn) {

	db, err := conn.GetDB(context.Background())
	if err != nil {
		log.Fatalf("error retrieving database conn DB to charge test data")
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failure to cast gormDB as sqlDB")
	}

	dbDriver, err := conn.GetDriver(context.Background())
	if err != nil {
		log.Fatalf("error retrieving database driver to charge test data")
	}

	switch dbDriver {
	case "postgres":
		driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			log.Fatalf("failure to get DB driver")
		}

		migrate, err := migrate.NewWithDatabaseInstance(
			postgresMigrationsPath,
			"postgres",
			driver,
		)
		if err != nil {
			log.Fatalf("failure to instantiate golang-migration: %s", err)
		}

		suite.migrate = migrate

		err = migrate.Up()
		if err != nil {
			log.Fatalf("failure to Up migrations: %s", err)
		}

		seedFileContent, err := os.ReadFile(postgresSeedPath)
		if err != nil {
			log.Fatalf("Erro ao ler o arquivo: %v", err)
		}

		if err = gormDB.Exec(string(seedFileContent)).Error; err != nil {
			log.Fatalf("failure to charge database: %s", err)
		}

	default:
		log.Fatalf("error connecting to database migrate to charge test data. Incorrect Driver: %s", dbDriver)
	}
}

func (suite *RepositoriesSuite) AccountRepositoryFindByUIDsuccess() {
	accountEntity, err := suite.AccountRepo.FindByUID(context.Background(), accountUID)
	assert.Equal(suite.T(), accountEntity.UID, accountUID)
	assert.NoError(suite.T(), err)

	suite.AccountEntity = accountEntity
}

func (suite *RepositoriesSuite) AccountRepositorySaveTransactionsSuccess() {
	transactionEntities := make(map[int]port.TransactionEntity)

	transactionEntities[1] = port.TransactionEntity{
		AccountID:    1,
		Amount:       decimal.NewFromFloat(100.00),
		MCC:          merchantCorrectMccToMap,
		MerchantName: merchantNameToMap,
		CategoryID:   merchantCategoryToMap,
	}

	err := suite.AccountRepo.SaveTransactions(context.Background(), transactionEntities)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoriesSuite) MerchantRepositoryFindByNameSuccess() {
	merchantEntity, err := suite.MerchantRepo.FindByName(context.Background(), merchantNameToMap)
	assert.Equal(suite.T(), merchantEntity.MCC, merchantCorrectMccToMap)
	assert.NoError(suite.T(), err)
}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}

func (suite *RepositoriesSuite) TestCases() {
	suite.T().Run("TestAccountRepositoryFindByUIDSuccess", func(t *testing.T) {
		suite.AccountRepositoryFindByUIDsuccess()
	})

	suite.T().Run("TestTransactionRepositorySaveSuccess", func(t *testing.T) {
		suite.AccountRepositorySaveTransactionsSuccess()
	})

	suite.T().Run("TestMerchantRepositoryFindByNameSuccess", func(t *testing.T) {
		suite.MerchantRepositoryFindByNameSuccess()
	})
}

func (suite *RepositoriesSuite) TearDownSuite() {
	err := suite.migrate.Down()
	if err != nil {
		log.Fatalf("failure to Down migrations: %s", err)
	}
}
