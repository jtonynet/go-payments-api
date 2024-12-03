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

	app := ctx.MustGet("app").(bootstrap.RESTApp)
	logger := app.Logger

	defer func() {
		elapsedTime := time.Now().UnixMilli() - timestamp
		debugLog(
			logger,
			fmt.Sprintf("Execution time: %d ms\n", elapsedTime),
		)
	}()

	var transactionRequest port.TransactionPaymentRequest
	if err := ctx.ShouldBindBodyWith(&transactionRequest, binding.JSON); err != nil {
		debugLog(
			logger,
			fmt.Sprintf("rejected: %s, error:%s ms\n", port.CODE_REJECTED_GENERIC, err.Error()),
		)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}

	validationErrors, ok := dtoIsValid(transactionRequest)
	if !ok {
		debugLog(logger, validationErrors)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}

	result, err := app.GRPCpayment.Execute(
		context.Background(),
		&pb.TransactionRequest{
			Account:     transactionRequest.AccountUID.String(),
			Mcc:         transactionRequest.MCC,
			Merchant:    transactionRequest.Merchant,
			TotalAmount: transactionRequest.TotalAmount.String(),
		},
	)

	if err != nil {
		debugLog(logger, err.Error())

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: port.CODE_REJECTED_GENERIC,
		})

		return
	}

	ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
		Code: result.Code,
	})
}

func debugLog(logger logger.Logger, msg string) {
	if logger != nil {
		logger.Debug(msg)
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
