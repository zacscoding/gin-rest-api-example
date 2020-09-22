package account

import (
	accountDB "gin-rest-api-example/internal/account/database"
	"gin-rest-api-example/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	accountDB accountDB.AccountDB
}

func (h *Handler) handleMe(ctx *gin.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("handle me..")

	acc, _ := CurrentUser(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"acc": acc,
	})
}

func (h *Handler) handleYou(ctx *gin.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("handle you..")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func NewHandler(accountDB accountDB.AccountDB) *Handler {
	return &Handler{
		accountDB: accountDB,
	}
}
