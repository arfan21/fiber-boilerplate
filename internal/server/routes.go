package server

import (
	userctrl "github.com/arfan21/fiber-boilerplate/internal/user/controller"
	userrepo "github.com/arfan21/fiber-boilerplate/internal/user/repository"
	usersvc "github.com/arfan21/fiber-boilerplate/internal/user/service"
	"github.com/arfan21/fiber-boilerplate/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userRepoRedis := userrepo.NewRedis(s.dbRedis)
	userSvc := usersvc.New(userRepo, userRepoRedis)
	userCtrl := userctrl.New(userSvc)

	s.RoutesCustomer(api, userCtrl)

}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/users")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
	usersV1.Post("/refresh-token", ctrl.RefreshToken)
	usersV1.Post("/logout", middleware.JWTAuth, ctrl.Logout)
}
