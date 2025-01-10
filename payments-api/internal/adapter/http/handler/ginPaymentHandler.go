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
	startTIme := time.Now().UnixMilli()
	code := port.CODE_REJECTED_GENERIC
	transactionUID := uuid.NewString()

	requestCtx := context.Background()
	requestCtx = context.WithValue(requestCtx, logger.CtxTransactionUIDKey, transactionUID)

	app := ctx.MustGet("app").(bootstrap.RESTApp)

	app.Logger.Info(
		requestCtx,
		"Transaction Initialized",
	)

	defer func() {
		elapsedTime := time.Now().UnixMilli() - startTIme
		requestCtx = context.WithValue(requestCtx, logger.CtxExecutionTimeKey, elapsedTime)
		requestCtx = context.WithValue(requestCtx, logger.CtxResponseCodeKey, code)
		app.Logger.Info(
			requestCtx,
			"Transaction Fineshed",
		)
	}()

	var transactionRequest port.TransactionPaymentRequest
	if err := ctx.ShouldBindBodyWith(&transactionRequest, binding.JSON); err != nil {
		app.Logger.Warn(
			context.Background(),
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
		app.Logger.Info(requestCtx, validationErrors)

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
		app.Logger.Warn(requestCtx, err.Error())

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
