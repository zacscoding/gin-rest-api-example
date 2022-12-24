package article

import (
	"context"
	"gin-rest-api-example/internal/account"
	articleDB "gin-rest-api-example/internal/article/database"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/internal/middleware"
	"gin-rest-api-example/internal/middleware/handler"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/validate"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type Handler struct {
	articleDB articleDB.ArticleDB
}

// saveArticle handles POST /v1/api/articles
func (h *Handler) saveArticle(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestBody struct {
			Article struct {
				Title string   `json:"title" binding:"required,min=5"`
				Body  string   `json:"body" binding:"required"`
				Tags  []string `json:"tagList" binding:"omitempty,dive,max=10"`
			} `json:"article"`
		}
		var body RequestBody
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Errorw("article.handler.register failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&body.Article, "json", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, "invalid article request in body", details)
		}

		// save article
		currentUser := account.MustCurrentUser(c)
		var tags []*model.Tag
		for _, tag := range body.Article.Tags {
			tags = append(tags, &model.Tag{Name: tag})
		}
		article := model.Article{
			Slug:     slug.Make(body.Article.Title),
			Title:    body.Article.Title,
			Body:     body.Article.Body,
			Author:   *currentUser,
			AuthorID: currentUser.ID,
			Tags:     tags,
		}
		err := h.articleDB.SaveArticle(c.Request.Context(), &article)
		if err != nil {
			if database.IsKeyConflictErr(err) {
				return handler.NewErrorResponse(http.StatusConflict, handler.DuplicateEntry, "duplicate article title", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusCreated, NewArticleResponse(&article))
	})
}

// articleBySlug handles GET /v1/api/articles/:slug
func (h *Handler) articleBySlug(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Slug string `uri:"slug"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.articleBySlug failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article request in uri", details)
		}

		// find
		article, err := h.articleDB.FindArticleBySlug(c.Request.Context(), uri.Slug)
		if err != nil {
			if database.IsRecordNotFoundErr(err) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found article", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewArticleResponse(article))
	})
}

// articles handles GET /v1/api/articles
func (h *Handler) articles(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		type QueryParameter struct {
			Tag    []string `form:"tag" binding:"omitempty,dive,max=10"`
			Author string   `form:"author" binding:"omitempty"`
			Limit  string   `form:"limit,default=5" binding:"numeric"`
			Offset string   `form:"offset,default=0" binding:"numeric"`
		}
		var query QueryParameter
		if err := c.ShouldBindQuery(&query); err != nil {
			logger.Errorw("article.handler.articles failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&query, "form", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article request in query", details)
		}

		limit, err := strconv.ParseUint(query.Limit, 10, 64)
		if err != nil {
			limit = 5
		}
		offset, err := strconv.ParseUint(query.Offset, 10, 64)
		if err != nil {
			offset = 0
		}
		criteria := articleDB.IterateArticleCriteria{
			Tags:   query.Tag,
			Author: query.Author,
			Offset: uint(offset),
			Limit:  uint(limit),
		}
		articles, total, err := h.articleDB.FindArticles(c.Request.Context(), criteria)
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, NewArticlesResponse(articles, total))
	})
}

// deleteArticle handles DELETE /v1/api/articles/:slug
func (h *Handler) deleteArticle(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		// bind
		type RequestUri struct {
			Slug string `uri:"slug" binding:"required"`
		}
		var uri RequestUri
		if err := c.ShouldBindUri(&uri); err != nil {
			logger.Errorw("article.handler.deleteArticle failed to bind", "err", err)
			var details []*validate.ValidationErrDetail
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				details = validate.ValidationErrorDetails(&uri, "uri", vErrs)
			}
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidUriValue, "invalid article request in uri", details)
		}

		// delete article and comments with in transaction
		currentUser := account.MustCurrentUser(c)
		err := h.articleDB.RunInTx(c.Request.Context(), func(ctx context.Context) error {
			// delete a article
			if err := h.articleDB.DeleteArticleBySlug(ctx, currentUser.ID, uri.Slug); err != nil {
				return err
			}

			// delete article comments
			deleted, err := h.articleDB.DeleteComments(ctx, currentUser.ID, uri.Slug)
			if err != nil {
				return err
			}
			logger.Debugw("article.handler.deleteArticle success to delete a article", "comments", deleted)
			return nil
		})
		if err != nil {
			logger.Errorw("article.handler.deleteArticle failed to delete a article", "err", err)
			if database.IsRecordNotFoundErr(errors.Cause(err)) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, "not found article", nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		return handler.NewSuccessResponse(http.StatusOK, nil)
	})
}

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine, auth *jwt.GinJWTMiddleware) {
	v1 := r.Group("v1/api")
	v1.Use(middleware.RequestIDMiddleware(), middleware.TimeoutMiddleware(cfg.ServerConfig.WriteTimeout))

	articleV1 := v1.Group("articles")
	// anonymous
	articleV1.Use()
	{
		articleV1.GET(":slug", h.articleBySlug)
		articleV1.GET("", h.articles)
		articleV1.GET(":slug/comments", h.articleComments)
	}

	// auth required
	articleV1.Use(auth.MiddlewareFunc())
	{
		articleV1.POST("", h.saveArticle)
		articleV1.DELETE(":slug", h.deleteArticle)
		articleV1.POST(":slug/comments", h.saveComment)
		articleV1.DELETE(":slug/comments/:id", h.deleteComment)
	}
}

func NewHandler(articleDB articleDB.ArticleDB) *Handler {
	return &Handler{
		articleDB: articleDB,
	}
}
