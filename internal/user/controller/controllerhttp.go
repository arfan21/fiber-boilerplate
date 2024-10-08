package userctrl

import (
	"github.com/arfan21/fiber-boilerplate/internal/model"
	usersvc "github.com/arfan21/fiber-boilerplate/internal/user/service"
	"github.com/arfan21/fiber-boilerplate/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc *usersvc.Service
}

func New(svc *usersvc.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Register user
// @Description Register user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserRegisterRequest true "Payload user Register Request"
// @Success 201 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/register [post]
func (ctrl ControllerHTTP) Register(c *fiber.Ctx) error {
	var req model.UserRegisterRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	err = ctrl.svc.Register(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusCreated,
	})
}

// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "Payload user Login Request"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.UserLoginResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/login [post]
func (ctrl ControllerHTTP) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	res, err := ctrl.svc.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusOK,
		Data: res,
	})
}

// @Summary Refresh Token user
// @Description Refresh Token user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserRefreshTokenRequest true "Payload user Refresh Token Request"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.UserLoginResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/refresh-token [post]
func (ctrl ControllerHTTP) RefreshToken(c *fiber.Ctx) error {
	var req model.UserRefreshTokenRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	res, err := ctrl.svc.RefreshToken(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusOK,
		Data: res,
	})
}

// @Summary Logout user
// @Description Logout user
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.UserLogoutRequest true "Payload user Logout Request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/logout [post]
func (ctrl ControllerHTTP) Logout(c *fiber.Ctx) error {
	var req model.UserLogoutRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	err = ctrl.svc.Logout(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusOK,
	})
}
