package ginHandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/google/uuid"

	"github.com/jtonynet/go-payments-api/bootstrap"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/jtonynet/go-payments-api/internal/support/logger"

	pb "github.com/jtonynet/go-payments-api/internal/adapter/gRPC/pb"
)

// @Summary Payment Execute Transaction
// @Description Payment executes a transaction  based on the request body json data. The HTTP status is always 200. The transaction can be **approved** (code **00**), **rejected insufficient balance** (code **51**), or **rejected generally** (code **07**). [See more here](https://github.com/jtonynet/go-payments-api/tree/main?tab=readme-ov-file#about)
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body port.TransactionPaymentRequest true "Request body for Execute Transaction Payment"
// @Router /payment [post]
// @Success 200 {object} port.TransactionPaymentResponse
func PaymentExecution(ctx *gin.Context) {
	timestamp := time.Now().UnixMilli()
	code := port.CODE_REJECTED_GENERIC

	app := ctx.MustGet("app").(bootstrap.RESTApp)
	loggerInstance := app.Logger

	transactionUID := uuid.NewString()

	requestCtx := context.Background()
	requestCtx = context.WithValue(requestCtx, logger.CtxTransactionUIDKey, transactionUID)
	defer func() {
		elapsedTime := time.Now().UnixMilli() - timestamp
		requestCtx = context.WithValue(requestCtx, logger.CtxExecutionTimeKey, elapsedTime)
		requestCtx = context.WithValue(requestCtx, logger.CtxResponseCodeKey, code)
		debugLog(
			requestCtx,
			loggerInstance,
			"Transaction Fineshed",
		)
	}()

	var transactionRequest port.TransactionPaymentRequest
	if err := ctx.ShouldBindBodyWith(&transactionRequest, binding.JSON); err != nil {
		debugLog(
			context.Background(),
			loggerInstance,
			fmt.Sprintf("rejected: %s, error:%s ms\n", port.CODE_REJECTED_GENERIC, err.Error()),
		)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}
	accountUID := transactionRequest.AccountUID.String()
	requestCtx = context.WithValue(requestCtx, logger.CtxAccountUIDKey, accountUID)

	validationErrors, ok := dtoIsValid(transactionRequest)
	if !ok {
		debugLog(context.Background(), loggerInstance, validationErrors)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}

	result, err := app.GRPCpayment.Execute(
		context.Background(),
		&pb.TransactionRequest{
			Transaction: transactionUID,
			Account:     accountUID,
			Mcc:         transactionRequest.MCC,
			Merchant:    transactionRequest.Merchant,
			TotalAmount: transactionRequest.TotalAmount.String(),
		},
	)

	if err != nil {
		debugLog(context.Background(), loggerInstance, err.Error())

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}

	code = result.Code
	ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
		Code: code,
	})
}

func debugLog(ctx context.Context, logger logger.Logger, msg string) {
	if logger != nil {
		logger.Debug(ctx, msg)
	}
}

func validateUUID(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(uuid.UUID)
	return ok
}

func dtoIsValid(dto any) (string, bool) {
	validate := validator.New()

	validate.RegisterValidation("uuid", validateUUID)

	err := validate.Struct(dto)
	if err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
		return strings.Join(errMsgs, ", "), false
	}

	return "", true
}
