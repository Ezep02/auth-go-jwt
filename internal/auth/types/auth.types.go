package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Age      int64
	Email    string
	Password string
	Admin    bool
	Is_user  bool
}

type UserRequest struct {
	Email    string
	Password string
}

type UserResponse struct {
	gorm.Model
	Name    string
	Age     int64
	Email   string
	Admin   bool
	Is_user bool
}

type LoginUserResponse struct {
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	UserResponse
}

type RenewAccesTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccesTokenRes struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

type UserSignInResponse struct {
	gorm.Model
	Name     string
	Password string
	Age      int64
	Email    string
	Admin    bool
	Is_user  bool
}

type Session struct {
	ID           string
	UserEmail    string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
