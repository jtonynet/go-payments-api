package gRPC

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/protobuffer"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/jtonynet/go-payments-api/internal/core/service"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	protobuffer.UnimplementedPaymentServer
	hostAndPort    string
	PaymentService service.Payment
}

func NewPaymentServer(cfg config.GRPC, PaymentService service.Payment) (PaymentServer, error) {
	return PaymentServer{
		hostAndPort:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		PaymentService: PaymentService,
	}, nil
}

func (ps *PaymentServer) HandleRequests() {
	listener, err := net.Listen("tcp", ps.hostAndPort)
	if err != nil {
		log.Fatalf("cannot initiate gRPC listner: %v", err)
	}

	s := grpc.NewServer()
	protobuffer.RegisterPaymentServer(s, ps)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
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

	code, _ := ps.PaymentService.Execute(
		port.TransactionPaymentRequest{
			AccountUID:  accountUID,
			TotalAmount: totalAmount,
			MCC:         tr.Mcc,
			Merchant:    tr.Merchant,
		},
	)

	return &protobuffer.TransactionResponse{Code: code}, nil
}
