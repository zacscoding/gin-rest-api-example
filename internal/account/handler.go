package account

import (
	accountDB "gin-rest-api-example/internal/account/database"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/internal/middleware"
	"gin-rest-api-example/internal/middleware/handler"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/validate"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Handler struct {
	accountDB accountDB.AccountDB
}

// signUp handles POST /v1/api/users
func (h *Handler) signUp(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		type RequestBody struct {
			User struct {
				Username string `json:"username" binding:"required"`
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required,min=5"`
			} `json:"user"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Errorw("account.handler.signUp failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&body.User, "json", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, "invalid user request in body", details)
		}

		password, err := EncodePassword(body.User.Password)
		if err != nil {
			logger.Errorw("account.handler.signUp failed to encode password", "err", err)
			return handler.NewInternalErrorResponse(err)
		}
		acc := model.Account{
			Username: body.User.Username,
			Email:    body.User.Email,
			Password: password,
		}
		err = h.accountDB.Save(c.Request.Context(), &acc)
		if err != nil {
			if database.IsKeyConflictErr(err) {
				return handler.NewErrorResponse(http.StatusConflict, handler.DuplicateEntry, "duplicate email address", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusCreated, NewUserResponse(&acc))
	})
}

// currentUser handles GET /v1/api/user/me
func (h *Handler) currentUser(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		currentUser := MustCurrentUser(c)
		find, err := h.accountDB.FindByEmail(c.Request.Context(), currentUser.Email)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found current user", nil)
			}
			return &handler.Response{Err: err}
		}
		if find.Disabled {
			return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found current user", nil)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewUserResponse(find))
	})
}

// update handles PUT /v1/api/user
func (h *Handler) update(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		currentUser := MustCurrentUser(c)
		type RequestBody struct {
			User struct {
				Username string `json:"username" binding:"omitempty"`
				Password string `json:"password" binding:"omitempty,min=5"`
				Bio      string `json:"bio"`
				Image    string `json:"image"`
			} `json:"user"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Errorw("account.handler.update failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&body.User, "json", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, "invalid user request in body", details)
		}

		acc, err := h.accountDB.FindByEmail(c.Request.Context(), currentUser.Email)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found account", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		if body.User.Password != "" {
			password, err := EncodePassword(body.User.Password)
			if err != nil {
				logger.Errorw("account.handler.update failed to encode password", "err", err)
				return &handler.Response{Err: err}
			}
			acc.Password = password
		}
		if body.User.Username != "" {
			acc.Username = body.User.Username
		}
		if body.User.Bio != "" {
			acc.Bio = body.User.Bio
		}
		if body.User.Image != "" {
			acc.Image = body.User.Image
		}
		err = h.accountDB.Update(c.Request.Context(), currentUser.Email, acc)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				logger.Errorw("account.handler.update failed to update user because not found user", "err", err)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewUserResponse(acc))
	})
}

// RouteV1 routes user api given config and gin.Engine
func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine, auth *jwt.GinJWTMiddleware) {
	v1 := r.Group("v1/api")
	timeout := time.Duration(cfg.ServerConfig.WriteTimeoutSecs) * time.Second
	v1.Use(middleware.RequestIDMiddleware(), middleware.TimeoutMiddleware(timeout))
	// anonymous
	v1.Use()
	{
		v1.POST("users/login", auth.LoginHandler)
		v1.POST("users", h.signUp)
	}
	// auth required
	v1.Use(auth.MiddlewareFunc())
	{
		v1.GET("user/me", h.currentUser)
		v1.PUT("user", h.update)
	}
}

func NewHandler(accountDB accountDB.AccountDB) *Handler {
	return &Handler{
		accountDB: accountDB,
	}
}
