package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/arfan21/fiber-boilerplate/pkg/constant"
	"github.com/arfan21/fiber-boilerplate/pkg/logger"
	"github.com/arfan21/fiber-boilerplate/pkg/pkgutil"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	defer func() {
		logger.Log(ctx.UserContext()).Error().Msg(err.Error())
	}()

	defaultRes := pkgutil.HTTPResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	var errValidation *constant.ErrValidation
	if errors.As(err, &errValidation) {
		data := errValidation.Error()
		var messages []map[string]interface{}

		errJson := json.Unmarshal([]byte(data), &messages)
		PanicIfNeeded(errJson)

		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = "Bad Request"
		var errors []interface{}
		for _, message := range messages {
			errors = append(errors, message)
		}
		defaultRes.Errors = errors
	}

	var withCodeErr *constant.ErrWithCode
	if errors.As(err, &withCodeErr) {
		defaultRes.Code = http.StatusInternalServerError
		if withCodeErr.HTTPStatusCode > 0 {
			defaultRes.Code = withCodeErr.HTTPStatusCode
		}
		defaultRes.Message = http.StatusText(defaultRes.Code)
		if withCodeErr.Message != "" {
			defaultRes.Message = withCodeErr.Message
		}
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		defaultRes.Code = fiberError.Code
		defaultRes.Message = fiberError.Message
	}

	if errors.Is(err, pgx.ErrNoRows) {
		defaultRes.Code = fiber.StatusNotFound
		defaultRes.Message = "data not found"
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		defaultRes.Code = fiber.StatusUnprocessableEntity
		defaultRes.Message = http.StatusText(fiber.StatusUnprocessableEntity)

		defaultRes.Errors = []interface{}{
			map[string]interface{}{
				"field":   unmarshalTypeError.Field,
				"message": fmt.Sprintf("%s harus %s", unmarshalTypeError.Field, unmarshalTypeError.Type),
			},
		}
	}

	// handle error parse uuid
	if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("invalid UUID")) {
		defaultRes.Code = fiber.StatusBadRequest
		defaultRes.Message = constant.ErrInvalidUUID.Error()
	}

	if defaultRes.Code >= 500 {
		defaultRes.Message = http.StatusText(defaultRes.Code)
	}

	return ctx.Status(defaultRes.Code).JSON(defaultRes)
}
