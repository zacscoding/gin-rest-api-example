package account

import (
	"gin-rest-api-example/pkg/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func RouteV1(h *Handler, r *gin.Engine, auth *jwt.GinJWTMiddleware) {
	v1 := r.Group("v1/api")
	v1.Use(middleware.RequestIDMiddleware())
	// anonymous
	v1.Use()
	{
		v1.POST("users/login", auth.LoginHandler)
		v1.GET("/you", h.handleYou)
	}
	// auth required
	v1.Use(auth.MiddlewareFunc())
	{
		v1.GET("user/me", h.handleMe)
	}
}
