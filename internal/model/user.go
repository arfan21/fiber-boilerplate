package model

type UserRegisterRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type UserLoginResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	ExpiresInRefreshToken int    `json:"expires_in_refresh_token,omitempty"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserLogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
