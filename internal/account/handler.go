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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func NewHandler(accountDB accountDB.AccountDB) *Handler {
	return &Handler{
		accountDB: accountDB,
	}
}
