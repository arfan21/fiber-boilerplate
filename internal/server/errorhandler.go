package server

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

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	defer func() {
		logger.Log(ctx.UserContext()).Error().Msg(err.Error())
	}()
	var arrWithCodeErr constant.ErrsWithCode

	var withCodeErr constant.ErrWithCode
	if errors.As(err, &withCodeErr) {
		arrWithCodeErr = append(arrWithCodeErr, withCodeErr)
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		arrWithCodeErr = append(arrWithCodeErr, constant.ErrWithCode{
			Message:        fiberError.Message,
			HTTPStatusCode: fiberError.Code,
		})
	}

	if errors.Is(err, pgx.ErrNoRows) {
		arrWithCodeErr = append(arrWithCodeErr, constant.ErrWithCode{
			Message:        "data not found",
			HTTPStatusCode: fiber.StatusNotFound,
		})
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		arrWithCodeErr = append(arrWithCodeErr, constant.ErrWithCode{
			Message:        fmt.Sprintf("%s should %s", unmarshalTypeError.Field, unmarshalTypeError.Type),
			HTTPStatusCode: fiber.StatusUnprocessableEntity,
			Field:          unmarshalTypeError.Field,
		})
	}

	// handle error parse uuid
	if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("invalid UUID")) {
		arrWithCodeErr = append(arrWithCodeErr, constant.ErrWithCode{
			Message:        constant.ErrInvalidUUID.Error(),
			HTTPStatusCode: fiber.StatusBadRequest,
		})
	}

	defaultRes := pkgutil.HTTPResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	if defaultRes.Code >= 500 {
		defaultRes.Message = http.StatusText(defaultRes.Code)
	}

	if len(arrWithCodeErr) > 0 {
		defaultRes.Code = arrWithCodeErr[0].HTTPStatusCode
		defaultRes.Message = arrWithCodeErr[0].Error()

		if arrWithCodeErr[0].Field != "" {
			defaultRes.Errors = arrWithCodeErr
		}
	}

	return ctx.Status(defaultRes.Code).JSON(defaultRes)
}
