package gRPC

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/protobuffer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPaymentClient(cfg config.GRPC) (protobuffer.PaymentClient, error) {
	gRPCServerConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	PaymentClient := protobuffer.NewPaymentClient(gRPCServerConn)

	return PaymentClient, nil
}
