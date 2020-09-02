package account

import "github.com/gin-gonic/gin"

func RouteV1(h *Handler, r *gin.Engine) {
	v1 := r.Group("v1")
	v1.Use()
	{
		v1.GET("/me", h.handleMe)
	}
}
