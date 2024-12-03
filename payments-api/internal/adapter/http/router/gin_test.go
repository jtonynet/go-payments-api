package router

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"

	pb "github.com/jtonynet/go-payments-api/internal/adapter/gRPC/pb"
	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"
)

var (
	accountUID, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	amountFoodTransaction = decimal.NewFromFloat(100.10)
)

type PaymentServerFake struct {
	pb.UnimplementedPaymentServer
}

func NewPaymentClientFake() (pb.PaymentClient, error) {
	return &PaymentServerFake{}, nil
}

func (ps *PaymentServerFake) Execute(
	ctx context.Context,
	tr *pb.TransactionRequest,
	opts ...grpc.CallOption,
) (*pb.TransactionResponse, error) {

	totalAmount, err := decimal.NewFromString(tr.TotalAmount)
	if err != nil {
		return nil, err
	}

	code := "00"

	if totalAmount.GreaterThan(amountFoodTransaction) {
		code = "51"
	}

	return &pb.TransactionResponse{Code: code}, nil
}

type GinRouterSuite struct {
	suite.Suite

	router   *gin.Engine
	apiGroup *gin.RouterGroup
}

func (suite *GinRouterSuite) SetupSuite() {

	cfg, err := config.LoadConfig("./../../../../")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	app, err := bootstrap.NewRESTApp(cfg)
	if err != nil {
		log.Fatal("cannot initiate app: ", err)
	}

	gRPCpayment, _ := NewPaymentClientFake()
	app.GRPCpayment = gRPCpayment

	suite.router, suite.apiGroup = setupRouterAndGroup(cfg.API, *app)

	suite.apiGroup.POST("/payment", ginHandler.PaymentExecution)
}

func setupRouterAndGroup(cfg config.API, app bootstrap.RESTApp) (*gin.Engine, *gin.RouterGroup) {
	basePath := "/"

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(ginMiddleware.ConfigInject(cfg))
	router.Use(ginMiddleware.AppInject(app))

	return router, router.Group(basePath)
}

func TestGinRouterSuite(t *testing.T) {
	suite.Run(t, new(GinRouterSuite))
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionApproved() {
	codeApproved := "00" // domain.CODE_APPROVED

	transactionJSON := fmt.Sprintf(
		`{
	  			"account": "%s",
	  			"mcc": "5411",
	  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
	  			"totalAmount": %v
			}`,
		accountUID,
		amountFoodTransaction,
	)

	suite.paymentExecuteTransactionTest(transactionJSON, codeApproved)
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionRejectedInsufficientFunds() {
	codeRejectedInsufficientFunds := "51" // domain.CODE_REJECTED_INSUFICIENT_FUNDS

	transactionJSON := fmt.Sprintf(
		`{
	  			"account": "%s",
	  			"mcc": "5411",
	  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
	  			"totalAmount": 9999.99
			}`,
		accountUID,
	)

	suite.paymentExecuteTransactionTest(transactionJSON, codeRejectedInsufficientFunds)
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionRejectedInvalidAccountUID() {
	codeRejectedInsufficientFunds := "07" // domain.CODE_REJECTED_GENERIC

	transactionJSON :=
		`{
			    "account": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	  			"mcc": "5411",
	  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
	  			"totalAmount": 0.01
			}`

	suite.paymentExecuteTransactionTest(transactionJSON, codeRejectedInsufficientFunds)
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionRejectedInvalidMCC() {
	codeRejectedInsufficientFunds := "07" // domain.CODE_REJECTED_GENERIC

	transactionJSON := fmt.Sprintf(
		`{
  			"account": "%s",
  			"mcc": "0",
  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
  			"totalAmount": 0.01
		}`,
		accountUID,
	)

	suite.paymentExecuteTransactionTest(transactionJSON, codeRejectedInsufficientFunds)
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionRejectedInvalidMerchant() {
	codeRejectedInsufficientFunds := "07" // domain.CODE_REJECTED_GENERIC

	transactionJSON := fmt.Sprintf(
		`{
  			"account": "%s",
  			"mcc": "5411",
  			"merchant": "PA",
  			"totalAmount": 0.01
		}`,
		accountUID,
	)

	suite.paymentExecuteTransactionTest(transactionJSON, codeRejectedInsufficientFunds)
}

func (suite *GinRouterSuite) TestPaymentExecuteTransactionRejectedInvalidAmount() {
	codeRejectedInsufficientFunds := "07" // domain.CODE_REJECTED_GENERIC

	transactionJSON := fmt.Sprintf(
		`{
  			"account": "%s",
  			"mcc": "5411",
  			"merchant": "PADARIA DO ZE              SAO PAULO BR",
  			"totalAmount": 00.01
		}`,
		accountUID,
	)

	suite.paymentExecuteTransactionTest(transactionJSON, codeRejectedInsufficientFunds)
}

func (suite *GinRouterSuite) paymentExecuteTransactionTest(reqBody string, returnCode string) {
	reqPaymentExecution, err := http.NewRequest("POST", "/payment", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	respPaymentExecution := httptest.NewRecorder()
	suite.router.ServeHTTP(respPaymentExecution, reqPaymentExecution)
	assert.Equal(suite.T(), http.StatusOK, respPaymentExecution.Code)

	bodyRespPaymentExecution := respPaymentExecution.Body.String()
	assert.Equal(suite.T(), gjson.Get(bodyRespPaymentExecution, "code").String(), returnCode)

}
