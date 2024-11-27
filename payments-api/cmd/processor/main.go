package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/protobuffer"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

/*
TODO WIP:
I still don't know where to fit this in the hexagonal. I'm searching
*/
type PaymentServer struct {
	protobuffer.UnimplementedPaymentServer
	app *bootstrap.ProcessorApp
}

func (ps *PaymentServer) Execute(
	ctx context.Context,
	tr *protobuffer.TransactionRequest,
) (*protobuffer.TransactionResponse, error) {

	accountUID, err := uuid.Parse(tr.Account)
	if err != nil {
		return nil, err
	}

	totalAmount, err := decimal.NewFromString(tr.TotalAmount)
	if err != nil {
		return nil, err
	}

	code, _ := ps.app.PaymentService.Execute(
		port.TransactionPaymentRequest{
			AccountUID:  accountUID,
			TotalAmount: totalAmount,
			MCC:         tr.Mcc,
			Merchant:    tr.Merchant,
		},
	)

	return &protobuffer.TransactionResponse{Code: code}, nil
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	app, err := bootstrap.NewProcessorApp(cfg)
	if err != nil {
		log.Fatalf("cannot initiate app: %v", err)
	}

	/*
	 TODO: gRPC Repository? Service?
	 Research how to best model this
	*/
	hostAndPort := fmt.Sprintf("%s:%s", cfg.GRPC.ServerHost, cfg.GRPC.ServerPort)
	listener, err := net.Listen("tcp", hostAndPort)
	if err != nil {
		log.Fatalf("cannot initiate gRPC listner: %v", err)
	}

	println("OK!!!")

	s := grpc.NewServer()
	protobuffer.RegisterPaymentServer(s, &PaymentServer{
		app: app,
	})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
