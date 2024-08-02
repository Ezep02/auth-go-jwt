package handler

import "github.com/gin-gonic/gin"

func AuthRouter(router *gin.Engine, handler *AuthHandler) *gin.Engine {

	auth_r := router.Group("/auth")
	{
		auth_r.POST("/sign-in", handler.UserSignIn)
		auth_r.POST("/sign-up", handler.UserSignUp)
		auth_r.POST("/logout", handler.LogoutUser)
	}

	token_r := router.Group("/token")
	{
		token_r.POST("/renew", handler.RenewAccesToken)
		token_r.POST("/revoke/:id", handler.RevokeSession)
	}

	return router
}
