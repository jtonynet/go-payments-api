package gRPC

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/protobuffer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPaymentClient(cfg config.GRPC) (protobuffer.PaymentClient, error) {
	hostAndPort := fmt.Sprintf("%s:%s", cfg.ClientHost, cfg.ClientPort)

	gRPCServerConn, err := grpc.Dial(
		hostAndPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	PaymentClient := protobuffer.NewPaymentClient(gRPCServerConn)

	return PaymentClient, nil
}
