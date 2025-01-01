package gRPC

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/config"
	pb "github.com/jtonynet/go-payments-api/internal/adapter/gRPC/pb"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/jtonynet/go-payments-api/internal/core/service"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServer
	hostAndPort    string
	paymentService service.Payment
}

func NewPaymentServer(cfg config.GRPC, paymentService service.Payment) (PaymentServer, error) {
	return PaymentServer{
		hostAndPort:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		paymentService: paymentService,
	}, nil
}

func (ps *PaymentServer) HandleRequests() {
	listener, err := net.Listen("tcp", ps.hostAndPort)
	if err != nil {
		log.Fatalf("cannot initiate gRPC listner: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServer(s, ps)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (ps *PaymentServer) Execute(
	ctx context.Context,
	tr *pb.TransactionRequest,
) (*pb.TransactionResponse, error) {

	accountUID, err := uuid.Parse(tr.Account)
	if err != nil {
		return nil, err
	}

	transactionUID, err := uuid.Parse(tr.Transaction)
	if err != nil {
		return nil, err
	}

	totalAmount, err := decimal.NewFromString(tr.TotalAmount)
	if err != nil {
		return nil, err
	}

	code, _ := ps.paymentService.Execute(
		port.TransactionPaymentRequest{
			AccountUID:     accountUID,
			TransactionUID: transactionUID,
			TotalAmount:    totalAmount,
			MCC:            tr.Mcc,
			Merchant:       tr.Merchant,
		},
	)

	return &pb.TransactionResponse{Code: code}, nil
}
