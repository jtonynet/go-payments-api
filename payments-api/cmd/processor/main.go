package main

import (
	"log"

	"github.com/jtonynet/go-payments-api/config"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/internal/adapter/gRPC"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	app, err := bootstrap.NewProcessorApp(cfg)
	if err != nil {
		log.Fatalf("cannot initiate app: %v", err)
	}

	gRPCPaymentServer, err := gRPC.NewPaymentServer(cfg.GRPC, *app.PaymentService)
	if err != nil {
		log.Fatalf("cannot initiate gRPCPaymentServer: %v", err)
	}
	gRPCPaymentServer.HandleRequests()

}
