package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ezep02/microservicios/internal/auth/service"
	"github.com/ezep02/microservicios/internal/auth/token"
	"github.com/ezep02/microservicios/internal/auth/types"
	"github.com/ezep02/microservicios/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service    *service.AuthService
	TokenMaker *token.JWTMaker
}

func NewAuthHandler(AuthService *service.AuthService, secretKey string) *AuthHandler {
	return &AuthHandler{
		Service:    AuthService,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (Auth_H *AuthHandler) UserSignIn(ctx *gin.Context) {
	var user types.UserRequest

	// Extraer datos del cuerpo de la solicitud
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// login usuario
	loggedUser, err := Auth_H.Service.UserSignIn(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log user"})
		return
	}

	// chequear si la contraseña es correcta
	err = utils.CheckPassword(user.Password, loggedUser.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "contraseña incorrecta",
		})
		return
	}

	// se crea un JWT
	accessToken, accessClaims, err := Auth_H.TokenMaker.CreateToken(int64(loggedUser.ID), loggedUser.Email, loggedUser.Admin, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating token",
		})
		return
	}

	// refresh Token
	refreshToken, refreshClaims, err := Auth_H.TokenMaker.CreateToken(int64(loggedUser.ID), loggedUser.Email, loggedUser.Admin, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating refresh token",
		})
		return
	}

	session, err := Auth_H.Service.CreateSession(&types.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    loggedUser.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating session",
		})
		return
	}

	res := types.LoginUserResponse{
		AccessToken:           accessToken,
		SessionID:             session.ID,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		UserResponse: types.UserResponse{
			Model:   loggedUser.Model,
			Name:    loggedUser.Name,
			Age:     loggedUser.Age,
			Email:   loggedUser.Email,
			Admin:   loggedUser.Admin,
			Is_user: loggedUser.Is_user,
		},
	}

	// Responder con el usuario creado
	ctx.JSON(http.StatusOK, res)
}

func (Auth_H *AuthHandler) UserSignUp(ctx *gin.Context) {
	var body types.User

	// Extraer datos del cuerpo de la solicitud
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	user := types.User{
		Model:    body.Model,
		Name:     body.Name,
		Age:      body.Age,
		Email:    body.Email,
		Password: hashedPassword,
		Admin:    body.Admin,
		Is_user:  body.Is_user,
	}

	// Agregar el nuevo usuario
	createdUser, err := Auth_H.Service.UserSignUp(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Responder con el usuario creado
	ctx.JSON(http.StatusOK, createdUser)
}

func (Auth_H *AuthHandler) LogoutUser(ctx *gin.Context) {
	var body struct {
		id string
	}

	// Extraer datos del cuerpo de la solicitud
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid logout id"})
		return
	}

	err := Auth_H.Service.DeleteSession(body.id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error deleting session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success logout",
	})
}

func (Auth_H *AuthHandler) RenewAccesToken(ctx *gin.Context) {

	var req types.RenewAccesTokenReq

	// Extraer datos del cuerpo de la solicitud
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid logout id"})
		return
	}

	refreshClaims, err := Auth_H.TokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := Auth_H.Service.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting the session"})
		return
	}

	if session.IsRevoked {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Session revoked"})
		return
	}

	if session.UserEmail != refreshClaims.Email {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session"})
		return
	}

	accessToken, accessClaim, err := Auth_H.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 24*time.Hour)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}

	res := types.RenewAccesTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaim.RegisteredClaims.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, res)
}

func (Auth_H *AuthHandler) RevokeSession(ctx *gin.Context) {

	id := ctx.Param("id")
	fmt.Println("id:", id)
	err := Auth_H.Service.RevokeSession(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "session revoked",
	})
}
