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
	log, err := initializeLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	gRPCPaymentClient, err := gRPC.NewPaymentClient(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC Client: %w", err)
	}

	return &RESTApp{
		Logger:      log,
		GRPCpayment: gRPCPaymentClient,
	}, nil
}

func NewProcessorApp(cfg *config.Config) (*ProcessorApp, error) {
	// Setting Value Objects
	timeoutSLA := port.TimeoutSLA(time.Duration(cfg.API.TimeoutSLA) * time.Millisecond)

	// Initialize supports
	log, err := initializeLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize adapters
	pubSubClient, err := initializePubSub(cfg.PubSub, log)
	if err != nil {
		return nil, err
	}

	lockClient, err := initializeDatabaseInMemory(cfg.Lock.ToInMemoryDatabase(), "Lock", log)
	if err != nil {
		return nil, err
	}

	cacheClient, err := initializeDatabaseInMemory(cfg.Cache.ToInMemoryDatabase(), "Cache", log)
	if err != nil {
		return nil, err
	}

	dbConn, err := initializeDatabase(cfg.Database, log)
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

	memoryLockRepo, err := repository.NewMemoryLock(lockClient, pubSubClient, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize memory lock repository: %w", err)
	}

	// Initialize services
	paymentService := service.NewPayment(
		timeoutSLA,
		allRepos.Account,
		cachedMerchantRepo,
		memoryLockRepo,
		log,
	)

	return &ProcessorApp{
		Logger:         log,
		PaymentService: paymentService,
	}, nil
}

func initializeLogger(cfg config.Logger) (logger.Logger, error) {
	log, err := logger.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	log.Debug(context.Background(), "Logger initialized successfully")
	return log, nil
}

func initializePubSub(cfg config.PubSub, log logger.Logger) (pubSub.PubSub, error) {
	pubsub, err := pubSub.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pub/sub client: %w", err)
	}

	log.Debug(context.Background(), "Pub/Sub client initialized successfully")
	return pubsub, nil
}

func initializeDatabaseInMemory(
	cfg config.InMemoryDatabase,
	componentName string,
	log logger.Logger,
) (database.InMemory, error) {
	conn, err := database.NewInMemory(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize %s client: %w", componentName, err)
	}

	if err := conn.Readiness(context.Background()); err != nil {
		return nil, fmt.Errorf("%s client is not ready: %w", componentName, err)
	}

	log.Debug(context.Background(), fmt.Sprintf("%s client initialized successfully", componentName))
	return conn, nil
}

func initializeDatabase(cfg config.Database, log logger.Logger) (database.Conn, error) {
	conn, err := database.NewConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	if err := conn.Readiness(context.Background()); err != nil {
		return nil, fmt.Errorf("database connection is not ready: %w", err)
	}

	log.Debug(context.Background(), "Database connection initialized successfully")
	return conn, nil
}
