package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-rest-api-example/internal/article/model"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"time"
)

var (
	dComment = model.Comment{
		ID:        1,
		Body:      "comment1",
		Slug:      dArticle.Slug,
		Author:    dUser,
		AuthorID:  dUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

func (s *HandlerSuite) TestSaveComment() {
	// given
	s.db.On("SaveComment", mock.Anything, dComment.Slug, mock.MatchedBy(func(c *model.Comment) bool {
		c.ID = dComment.ID
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		return true
	})).Return(nil)

	// when
	requestBody := map[string]interface{}{
		"comment": map[string]interface{}{
			"body": dComment.Body,
		},
	}
	b, _ := json.Marshal(&requestBody)
	url := fmt.Sprintf("/v1/api/articles/%s/comments", dComment.Slug)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Add("Authorization", "Bearer "+s.getBearerToken())

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "SaveComment", mock.Anything, dComment.Slug, mock.Anything)
	// 2) status code
	s.Equal(http.StatusCreated, res.Code)
	// 3) response
	s.assertCommentResponse(&dComment, gjson.Parse(res.Body.String()).Get("comment"))
}

func (s *HandlerSuite) TestArticleComments() {
	// given
	s.db.On("FindComments", mock.Anything, dComment.Slug).Return([]*model.Comment{&dComment}, nil)

	// when
	url := fmt.Sprintf("/v1/api/articles/%s/comments", dComment.Slug)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "FindComments", mock.Anything, dComment.Slug)
	// 2) status code
	s.Equal(http.StatusOK, res.Code)
	// 3) response
	results := gjson.Parse(res.Body.String()).Get("comments").Array()
	s.Equal(1, len(results))
	s.assertCommentResponse(&dComment, results[0])
}

func (s *HandlerSuite) TestDeleteComment() {
	// given
	s.db.On("DeleteCommentById", mock.Anything, dComment.Author.ID, dComment.Slug, dComment.ID).Return(nil)

	// when
	url := fmt.Sprintf("/v1/api/articles/%s/comments/%d", dComment.Slug, dComment.ID)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+s.getBearerToken())

	s.r.ServeHTTP(res, req)

	// then
	// 1) method called
	s.db.AssertCalled(s.T(), "DeleteCommentById", mock.Anything, dComment.Author.ID, dComment.Slug, dComment.ID)
	// 2) status code
	s.Equal(http.StatusOK, res.Code)
	// 3) response
	s.Empty(res.Body.Bytes())
}

func (s *HandlerSuite) assertCommentResponse(comment *model.Comment, result gjson.Result) {
	s.Equal(int64(comment.ID), result.Get("id").Int())
	s.True(result.Get("createdAt").Exists())
	s.True(result.Get("updatedAt").Exists())
	s.Equal(comment.Body, result.Get("body").String())
	s.Equal(comment.Author.Username, result.Get("author.username").String())
	s.Equal(comment.Author.Bio, result.Get("author.bio").String())
	s.Equal(comment.Author.Image, result.Get("author.image").String())
}
