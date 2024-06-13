package user

import (
	"context"

	"github.com/arfan21/fiber-boilerplate/internal/model"
)

type Service interface {
	Register(ctx context.Context, req model.UserRegisterRequest) (err error)
	Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error)
	RefreshToken(ctx context.Context, req model.UserRefreshTokenRequest) (res model.UserLoginResponse, err error)
	Logout(ctx context.Context, req model.UserLogoutRequest) (err error)
}
