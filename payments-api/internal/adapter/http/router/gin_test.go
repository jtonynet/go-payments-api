package router

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/adapter/database"

	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"
)

var (
	accountUID, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	balanceFoodUID, _ = uuid.Parse("bf64460c-efe0-4ffb-bd1f-54b136c2b2ac")
	balanceMealUID, _ = uuid.Parse("19475f4b-ee4c-4bce-add1-df0db5908201")
	balanceCashUID, _ = uuid.Parse("389e9316-ce28-478e-b14e-f971812de22d")

	balanceFoodAmount = decimal.NewFromFloat(205.11)
	balanceMealAmount = decimal.NewFromFloat(110.22)
	balanceCashAmount = decimal.NewFromFloat(115.33)

	amountFoodTransaction = decimal.NewFromFloat(100.10)
)

type GinRouterSuite struct {
	suite.Suite

	router   *gin.Engine
	apiGroup *gin.RouterGroup
}

func (suite *GinRouterSuite) SetupSuite() {

	/*
	   TODO:
	   Avoid concurrent testing with repositories
	   Find a better way to handle this
	*/
	time.Sleep(1 * time.Second)

	cfg, err := config.LoadConfig("./../../../../")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	app, err := bootstrap.NewRESTApp(cfg)
	if err != nil {
		log.Fatal("cannot initiate app: ", err)
	}

	conn, err := database.NewConn(cfg.Database)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	suite.loadDBtestData(conn)
	suite.router, suite.apiGroup = setupRouterAndGroup(cfg.API, *app)

	suite.apiGroup.POST("/payment", ginHandler.PaymentExecution)
}

func (suite *GinRouterSuite) loadDBtestData(conn database.Conn) {
	strategy, err := conn.GetStrategy(context.Background())
	if err != nil {
		log.Fatalf("error retrieving database strategy: %v", err)
	}

	switch strategy {
	case "gorm":
		db, err := conn.GetDB(context.Background())
		if err != nil {
			log.Fatalf("error retrieving database conn DB: %v", err)
		}

		dbGorm, ok := db.(*gorm.DB)
		if !ok {
			log.Fatalf("failure to cast conn.GetDB() as gorm.DB")
		}

		dbGorm.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE")
		insertCategoryQuery := `
			INSERT INTO categories (uid, name, priority, created_at, updated_at)
			VALUES
				('5681b4b5-6176-498a-a856-8932f79c05cc', 'FOOD', 1, NOW(), NOW()),
				('7bcfcd2a-2fde-4564-916b-92410e794272', 'MEAL', 2, NOW(), NOW()),
				('056de185-bff0-4c4a-93fa-7245f9e72b67', 'CASH', 3, NOW(), NOW())
		`
		dbGorm.Exec(insertCategoryQuery)

		dbGorm.Exec("TRUNCATE TABLE mccs RESTART IDENTITY CASCADE")
		insertMCCQuery := `
			INSERT INTO mccs (uid, mcc, category_id, created_at, updated_at)
			VALUES
				('11f0c06e-0dff-4643-86bf-998d11e9374f', '5411', 1, NOW(), NOW()),
				('fe5a4c17-a7cd-4072-a793-e99e2642e21a', '5412', 1, NOW(), NOW()),
				('5268ec2b-aa14-4d55-906a-13c91d89826c', '5811', 2, NOW(), NOW()),
				('6179e57c-e630-4e2f-a5db-d153e0cdb9a9', '5812', 2, NOW(), NOW())
		`
		dbGorm.Exec(insertMCCQuery)

		dbGorm.Exec("TRUNCATE TABLE merchants RESTART IDENTITY CASCADE")
		insertMerchantQuery := `
			INSERT INTO merchants (uid, name, mcc_id, created_at, updated_at)
			VALUES
				('95abe1ff-6f67-4a17-a4eb-d4842e324f1f', 'UBER EATS                   SAO PAULO BR', 2, NOW(), NOW()),
				('a53c6a52-8a18-4e7d-8827-7f612233c7ec', 'PAG*JoseDaSilva          RIO DE JANEI BR', 4, NOW(), NOW())`
		dbGorm.Exec(insertMerchantQuery)

		dbGorm.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")
		insertAccountQuery := fmt.Sprintf(
			`INSERT INTO accounts (uid, name, created_at, updated_at) 
			 VALUES('%s', 'Jonh Doe', NOW(), NOW())`,
			accountUID)
		dbGorm.Exec(insertAccountQuery)

		dbGorm.Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE")
		insertTransactionQuery := fmt.Sprintf(`
			INSERT INTO transactions (account_id, amount, category_id, created_at, updated_at)
			VALUES
				(1, %v, 1, NOW(), NOW()),
				(1, %v, 2, NOW(), NOW()),
				(1, %v, 3, NOW(), NOW())`,
			balanceFoodAmount,
			balanceMealAmount,
			balanceCashAmount,
		)
		dbGorm.Exec(insertTransactionQuery)

		dbGorm.Exec("TRUNCATE TABLE account_categories RESTART IDENTITY CASCADE")
		insertAccountCategoriesQuery := `
			INSERT INTO account_categories (account_id, category_id, created_at, updated_at)
			VALUES
				(1, 1, NOW(), NOW()),
				(1, 2, NOW(), NOW()),
				(1, 3, NOW(), NOW())`
		dbGorm.Exec(insertAccountCategoriesQuery)

	default:
		log.Fatalf("error connecting to database migrate to charge test data")
	}
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
