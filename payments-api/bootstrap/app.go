package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/config"

	"github.com/jtonynet/go-payments-api/internal/support/logger"

	"github.com/jtonynet/go-payments-api/internal/adapter/database"
	"github.com/jtonynet/go-payments-api/internal/adapter/gRPC"
	pb "github.com/jtonynet/go-payments-api/internal/adapter/gRPC/pb"
	"github.com/jtonynet/go-payments-api/internal/adapter/pubSub"
	"github.com/jtonynet/go-payments-api/internal/adapter/repository"

	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/jtonynet/go-payments-api/internal/core/service"
)

type RESTApp struct {
	Logger logger.Logger

	GRPCpayment pb.PaymentClient
}

type ProcessorApp struct {
	Logger logger.Logger

	PaymentService *service.Payment
}

func NewRESTApp(cfg *config.Config) (*RESTApp, error) {
	logger, err := initializeLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	gRPCPaymentClient, err := gRPC.NewPaymentClient(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC Client: %w", err)
	}

	return &RESTApp{
		Logger:      logger,
		GRPCpayment: gRPCPaymentClient,
	}, nil
}

func NewProcessorApp(cfg *config.Config) (*ProcessorApp, error) {
	timeoutSLA := port.TimeoutSLA(time.Duration(cfg.API.TimeoutSLA) * time.Millisecond)

	logger, err := initializeLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize adapters
	pubSubClient, err := initializePubSub(cfg.PubSub, logger)
	if err != nil {
		return nil, err
	}

	lockClient, err := initializeDatabaseInMemory(cfg.Lock.ToInMemoryDatabase(), "lock", logger)
	if err != nil {
		return nil, err
	}

	cacheClient, err := initializeDatabaseInMemory(cfg.Cache.ToInMemoryDatabase(), "cache", logger)
	if err != nil {
		return nil, err
	}

	dbConn, err := initializeDatabase(cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	allRepos, err := repository.GetAll(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	cachedMerchantRepo, err := repository.NewCachedMerchant(cacheClient, allRepos.Merchant)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cached merchant repository: %w", err)
	}

	memoryLockRepo, err := repository.NewMemoryLock(lockClient, pubSubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize memory lock repository: %w", err)
	}

	// Initialize services
	paymentService := service.NewPayment(
		timeoutSLA,
		allRepos.Account,
		cachedMerchantRepo,
		memoryLockRepo,
		logger,
	)

	return &ProcessorApp{
		Logger:         logger,
		PaymentService: paymentService,
	}, nil
}

func initializeLogger(cfg config.Logger) (logger.Logger, error) {
	logger, err := logger.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Debug("Logger initialized successfully")
	return logger, nil
}

func initializePubSub(cfg config.PubSub, logger logger.Logger) (pubSub.PubSub, error) {
	pubsub, err := pubSub.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pub/sub client: %w", err)
	}

	logger.Debug("Pub/Sub client initialized successfully")
	return pubsub, nil
}

func initializeDatabaseInMemory(
	cfg config.InMemoryDatabase,
	componentName string,
	logger logger.Logger,
) (database.InMemory, error) {
	conn, err := database.NewInMemory(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s client: %w", componentName, err)
	}

	if err := conn.Readiness(context.Background()); err != nil {
		return nil, fmt.Errorf("%s client is not ready: %w", componentName, err)
	}

	logger.Debug(fmt.Sprintf("%s client initialized successfully", componentName))
	return conn, nil
}

func initializeDatabase(cfg config.Database, logger logger.Logger) (database.Conn, error) {
	conn, err := database.NewConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	if err := conn.Readiness(context.Background()); err != nil {
		return nil, fmt.Errorf("database connection is not ready: %w", err)
	}

	logger.Debug("Database connection initialized successfully")
	return conn, nil
}
