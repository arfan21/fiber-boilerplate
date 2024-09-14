package usersvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/fiber-boilerplate/config"
	"github.com/arfan21/fiber-boilerplate/internal/entity"
	"github.com/arfan21/fiber-boilerplate/internal/model"
	userrepo "github.com/arfan21/fiber-boilerplate/internal/user/repository"
	"github.com/arfan21/fiber-boilerplate/pkg/constant"
	"github.com/arfan21/fiber-boilerplate/pkg/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *userrepo.Repository
	repoRedis *userrepo.RepositoryRedis
}

func New(repo *userrepo.Repository, repoRedis *userrepo.RepositoryRedis) *Service {
	return &Service{repo: repo, repoRedis: repoRedis}
}

func (s Service) Register(ctx context.Context, req model.UserRegisterRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to validate request: %w", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to hash password: %w", err)
		return
	}

	fmt.Println("hashedPassword", string(hashedPassword))

	data := entity.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to register user: %w", err)
		return
	}

	return
}

func (s Service) Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to validate request: %w", err)
		return
	}

	data, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to get user by email: %w", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = constant.ErrEmailOrPasswordInvalid
		}
		err = fmt.Errorf("user.service.Login: failed to compare password: %w", err)
		return
	}

	accessTokenExpire := time.Duration(config.Get().JWT.AccessTokenExpireIn) * time.Second

	accessToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Email,
		config.Get().JWT.AccessTokenSecret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to create access token: %w", err)
		return
	}

	refreshTokenExpire := time.Duration(config.Get().JWT.RefreshTokenExpireIn) * time.Second

	refreshToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Email,
		config.Get().JWT.RefreshTokenSecret,
		refreshTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to create refresh token: %w", err)
		return
	}

	err = s.repoRedis.SetRefreshToken(ctx, refreshToken, refreshTokenExpire, entity.UserRefreshToken{
		Email: data.Email,
		ID:    data.ID,
	})
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to set refresh token: %w", err)
		return
	}

	res = model.UserLoginResponse{
		AccessToken:           accessToken,
		ExpiresIn:             int(accessTokenExpire.Seconds()),
		TokenType:             "Bearer",
		RefreshToken:          refreshToken,
		ExpiresInRefreshToken: int(refreshTokenExpire.Seconds()),
	}

	return
}

func (s Service) CreateJWTWithExpiry(id, email, secret string, expiry time.Duration) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Get().Service.Name,
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("usecase: failed to create jwt token: %w", err)
		return
	}

	return
}

func (s Service) RefreshToken(ctx context.Context, req model.UserRefreshTokenRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.RefreshToken: failed to validate request: %w", err)
		return
	}

	payload, err := s.repoRedis.IsRefreshTokenExist(ctx, req.RefreshToken)
	if err != nil {
		err = fmt.Errorf("user.service.RefreshToken: failed to check refresh token: %w", err)
		return
	}

	accessTokenExpire := time.Duration(config.Get().JWT.AccessTokenExpireIn) * time.Second
	accessToken, err := s.CreateJWTWithExpiry(
		payload.ID.String(),
		payload.Email,
		config.Get().JWT.AccessTokenSecret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.RefreshToken: failed to create access token: %w", err)
		return
	}

	res = model.UserLoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(accessTokenExpire.Seconds()),
		TokenType:   "Bearer",
	}

	return
}

func (s Service) Logout(ctx context.Context, req model.UserLogoutRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Logout: failed to validate request: %w", err)
		return
	}

	err = s.repoRedis.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		err = fmt.Errorf("user.service.Logout: failed to delete refresh token: %w", err)
		return
	}

	return
}
