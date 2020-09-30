package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-rest-api-example/internal/account"
	accountDBMock "gin-rest-api-example/internal/account/database/mocks"
	accountModel "gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/article/database"
	articleDBMock "gin-rest-api-example/internal/article/database/mocks"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	dUser = accountModel.Account{
		ID:       1,
		Username: "user1",
		Email:    "user1@gmail.com",
		Password: "$2a$10$lsYsLv8nGPM0.R.ft4sgpe3OP7..KL3ZJqqhSVCKTEnSCMUztoUcW",
		Bio:      "I am working!",
	}
	dUserRawPass = "user1"

	dArticle = model.Article{
		ID:        1,
		Slug:      "how-to-train-your-dragon",
		Title:     "How to train your dragon",
		Body:      "You have to believe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Author:    dUser,
		AuthorID:  dUser.ID,
		Tags: []*model.Tag{
			{
				ID:        3,
				Name:      "reactjs",
				CreatedAt: time.Now(),
			}, {
				ID:        2,
				Name:      "angularjs",
				CreatedAt: time.Now(),
			}, {
				ID:        1,
				Name:      "dragons",
				CreatedAt: time.Now(),
			},
		},
	}
	dArticleTags = []string{"reactjs", "angularjs", "dragons"}
)

type HandlerSuite struct {
	suite.Suite
	r         *gin.Engine
	handler   *Handler
	db        *articleDBMock.ArticleDB
	accountDB *accountDBMock.AccountDB
}

func (s *HandlerSuite) SetupSuite() {
	logging.SetLevel(zapcore.FatalLevel)
}

func (s *HandlerSuite) SetupTest() {
	cfg, err := config.Load("")
	s.NoError(err)

	s.db = &articleDBMock.ArticleDB{}
	s.handler = NewHandler(s.db)
	s.accountDB = &accountDBMock.AccountDB{}
	s.accountDB.On("FindByEmail", mock.Anything, mock.MatchedBy(func(email string) bool {
		return email == dUser.Email
	})).Return(&dUser, nil)

	jwtMiddleware, err := account.NewAuthMiddleware(cfg, s.accountDB)
	s.NoError(err)

	gin.SetMode(gin.TestMode)
	s.r = gin.Default()

	RouteV1(cfg, s.handler, s.r, jwtMiddleware)

	accountHandler := account.NewHandler(s.accountDB)
	account.RouteV1(cfg, accountHandler, s.r, jwtMiddleware)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

// TODO : test failures

func (s *HandlerSuite) TestSaveArticle() {
	// given
	s.db.On("SaveArticle", mock.Anything, mock.Anything).Return(nil)

	// when
	requestBody := map[string]interface{}{
		"article": map[string]interface{}{
			"title":   dArticle.Title,
			"body":    dArticle.Body,
			"tagList": dArticleTags,
		},
	}
	b, _ := json.Marshal(&requestBody)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/api/articles", bytes.NewBuffer(b))
	req.Header.Add("Authorization", "Bearer "+s.getBearerToken())

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	articleMatcher := articleMatcher(dArticle.Title, dArticle.Body, dArticleTags, &dUser)
	s.db.AssertCalled(s.T(), "SaveArticle", mock.Anything, mock.MatchedBy(articleMatcher))
	// 2) status code
	s.Equal(http.StatusCreated, res.Code)
	// 3) response
	jsonVal := res.Body.String()
	s.assertArticleResponse(&dArticle, gjson.Parse(jsonVal).Get("article"))
}

func (s *HandlerSuite) TestArticleBySlug() {
	// given
	s.db.On("FindArticleBySlug", mock.Anything, dArticle.Slug).Return(&dArticle, nil)

	// when
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/api/articles/"+dArticle.Slug, nil)

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "FindArticleBySlug", mock.Anything, dArticle.Slug)
	// 2) status code
	s.Equal(http.StatusOK, res.Code)
	// 3) body
	jsonVal := res.Body.String()
	s.assertArticleResponse(&dArticle, gjson.Parse(jsonVal).Get("article"))
}

func (s *HandlerSuite) TestArticles() {
	criteria := database.IterateArticleCriteria{
		Tags:   []string{dArticleTags[0]},
		Author: dArticle.Author.Username,
		Offset: 0,
		Limit:  5,
	}
	s.db.On("FindArticles", mock.Anything, criteria).Return([]*model.Article{&dArticle}, int64(1), nil)

	// when
	url := fmt.Sprintf("/v1/api/articles?tag=%s&author=%s&offset=%d&limit=%d", criteria.Tags[0],
		criteria.Author, criteria.Offset, criteria.Limit)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "FindArticles", mock.Anything, criteria)
	// 2) status code
	s.Equal(http.StatusOK, res.Code)
	// 3) body
	jsonVal := res.Body.String()
	result := gjson.Parse(jsonVal)
	s.Equal(int64(1), result.Get("articlesCount").Int())

	articlesResult := result.Get("articles").Array()
	s.Equal(1, len(articlesResult))
	s.assertArticleResponse(&dArticle, articlesResult[0])
}

func (s *HandlerSuite) TestDeleteArticle() {
	// given
	s.db.On("RunInTx", mock.Anything, mock.Anything).Return(nil)

	// when
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/api/articles/"+dArticle.Slug, nil)
	req.Header.Add("Authorization", "Bearer "+s.getBearerToken())

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "RunInTx", mock.Anything, mock.Anything)
	// 2) status code
	s.Equal(http.StatusOK, res.Code)
	// 3) body
	s.Empty(res.Body.Bytes())
}

func (s *HandlerSuite) assertArticleResponse(article *model.Article, result gjson.Result) {
	s.Equal(slug.Make(article.Title), result.Get("slug").String())
	s.Equal(article.Title, result.Get("title").String())
	s.Equal(article.Body, result.Get("body").String())

	var tagVals []string
	for _, tag := range article.Tags {
		tagVals = append(tagVals, tag.Name)
	}
	for _, result := range result.Get("tagList").Array() {
		s.Contains(tagVals, result.String())
	}
	s.True(result.Get("createdAt").Exists())
	s.True(result.Get("updatedAt").Exists())
	s.Equal(article.Author.Username, result.Get("author.username").String())
	s.Equal(article.Author.Bio, result.Get("author.bio").String())
	s.Equal(article.Author.Image, result.Get("author.image").String())
}

func (s *HandlerSuite) getBearerToken() string {
	body := map[string]interface{}{
		"user": map[string]interface{}{
			"email":    dUser.Email,
			"password": dUserRawPass,
		},
	}
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/api/users/login", bytes.NewBuffer(b))
	s.r.ServeHTTP(res, req)

	s.Equal(http.StatusOK, res.Code)
	return gjson.Get(res.Body.String(), "token").String()
}

func articleMatcher(title, body string, tags []string, author *accountModel.Account) func(a *model.Article) bool {
	return func(a *model.Article) bool {
		if a.Slug != slug.Make(title) || a.Title != title || a.Body != body {
			return false
		}
		if len(tags) != len(a.Tags) {
			return false
		}
		tagVals := strings.Join(tags, " ")
		for _, tag := range a.Tags {
			if !strings.Contains(tagVals, tag.Name) {
				return false
			}
		}
		if a.Author.Username != author.Username || a.Author.Bio != author.Bio || a.Author.Image != author.Image {
			return false
		}
		return true
	}
}
