package gRPC

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	pb "github.com/jtonynet/go-payments-api/internal/adapter/gRPC/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPaymentClient(cfg config.GRPC) (pb.PaymentClient, error) {
	hostAndPort := fmt.Sprintf("%s:%s", cfg.ClientHost, cfg.ClientPort)

	gRPCClientConn, err := grpc.Dial(
		hostAndPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	PaymentClient := pb.NewPaymentClient(gRPCClientConn)

	return PaymentClient, nil
}
