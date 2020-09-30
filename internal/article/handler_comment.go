package article

import (
	"gin-rest-api-example/internal/account"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/internal/middleware/handler"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

// saveComment handles POST /v1/api/articles/:slug/comments
func (h *Handler) saveComment(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Slug string `uri:"slug" binding:"required"`
		}
		type RequestBody struct {
			Comment struct {
				Body string `json:"body" binding:"required"`
			} `json:"comment" binding:"required"`
		}
		var (
			uri  RequestUri
			body RequestBody
		)
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.saveComment failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article comment request in uri", details)
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Errorw("article.handler.saveComment failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "json", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article comment request in body", details)
		}

		currentUser := account.MustCurrentUser(c)
		comment := model.Comment{
			Body:   body.Comment.Body,
			Author: *currentUser,
		}
		if err := h.articleDB.SaveComment(c, uri.Slug, &comment); err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found article", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusCreated, NewCommentResponse(&comment))
	})
}

// articleComment handles GET /v1/api/articles/:slug/comments
func (h *Handler) articleComments(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Slug string `uri:"slug"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.articleComments failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article comment request in uri", details)
		}

		comments, err := h.articleDB.FindComments(c.Request.Context(), uri.Slug)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found article", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewCommentsResponse(comments))
	})
}

func (h *Handler) deleteComment(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Slug string `uri:"slug"`
			ID   string `uri:"id" binding:"numeric"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.deleteComment failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article comment request in uri", details)
		}
		id, err := strconv.ParseUint(uri.ID, 10, 64)
		if err != nil {
			details := validate.NewValidationErrorDetails("id", "id must be greater than or equals to 0", uri.ID)
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article comment request in uri", details)
		}

		// delete
		currentUser := account.MustCurrentUser(c)
		if err := h.articleDB.DeleteCommentById(c.Request.Context(), currentUser.ID, uri.Slug, uint(id)); err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found article comment", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, nil)
	})
}
