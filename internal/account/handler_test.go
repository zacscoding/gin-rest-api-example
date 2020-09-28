package account

import (
	"bytes"
	"encoding/json"
	"gin-rest-api-example/internal/account/database/mocks"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type HandlerSuite struct {
	suite.Suite
	r       *gin.Engine
	handler *Handler
	db      *mocks.AccountDB
}

func (s *HandlerSuite) SetupTest() {
	cfg, err := config.Load("")
	s.NoError(err)

	s.db = &mocks.AccountDB{}
	s.handler = NewHandler(s.db)

	jwtMiddleware, err := NewAuthMiddleware(cfg, s.db)
	s.NoError(err)

	gin.SetMode(gin.TestMode)
	s.r = gin.Default()

	RouteV1(cfg, s.handler, s.r, jwtMiddleware)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) TestRegister() {
	// given
	username := "zaccoding"
	email := "zaccoding@gmail.com"
	password := "password123"

	matcher := func(acc *model.Account) bool {
		return acc.Username == username && acc.Email == email && MatchesPassword(acc.Password, password) == nil
	}
	s.db.On("Save", mock.Anything, mock.MatchedBy(matcher)).Return(nil)

	// when
	body := map[string]interface{}{
		"user": map[string]interface{}{
			"username": username,
			"email":    email,
			"password": password,
		},
	}
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/api/users", bytes.NewBuffer(b))

	s.r.ServeHTTP(res, req)

	// then
	s.db.AssertCalled(s.T(), "Save", mock.Anything, mock.MatchedBy(matcher))
	s.Equal(http.StatusCreated, res.Code)
	expected := `
		{
		  "user": {
			"username": "zaccoding",
			"email": "zaccoding@gmail.com",
			"bio": "",
			"image": ""
		  }
		}`
	s.JSONEq(expected, res.Body.String())
}

func (s *HandlerSuite) TestRegister_BadRequest() {
	username := "user1"
	email := "user@email.com"
	password := "password123"
	cases := []struct {
		Username string
		Email    string
		Password string
		Expected string
	}{
		{
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "message": "[InvalidBodyValue] invalid user request in body",
			  "errors": [
				{
				  "field": "username",
				  "value": "",
				  "message": "required username"
				},
				{
				  "field": "email",
				  "value": "",
				  "message": "required email"
				},
				{
				  "field": "password",
				  "value": "",
				  "message": "required password"
				}
			  ]
			}`,
		},
		// username
		{
			Email:    email,
			Password: password,
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "message": "[InvalidBodyValue] invalid user request in body",
			  "errors": [
				{
				  "field": "username",
				  "value": "",
				  "message": "required username"
				}
			  ]
			}`,
		},
		// email
		{
			Username: username,
			Password: password,
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "message": "[InvalidBodyValue] invalid user request in body",
			  "errors": [
				{
				  "field": "email",
				  "value": "",
				  "message": "required email"
				}
			  ]
			}`,
		}, {
			Username: username,
			Email:    "invalid-email-format",
			Password: password,
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "message": "[InvalidBodyValue] invalid user request in body",
			  "errors": [
				{
				  "field": "email",
				  "value": "invalid-email-format",
				  "message": "required email format"
				}
			  ]
			}`,
		},
		// password
		{
			Username: username,
			Email:    email,
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "errors": [
				{
				  "field": "password",
				  "value": "",
				  "message": "required password"
				}
			  ],
			  "message": "[InvalidBodyValue] invalid user request in body"
			}`,
		}, {
			Username: username,
			Email:    email,
			Password: "1",
			Expected: `
			{
			  "code": "InvalidBodyValue",
			  "errors": [
				{
				  "field": "password",
				  "value": "1",
				  "message": "password required at least 5 length"
				}
			  ],
			  "message": "[InvalidBodyValue] invalid user request in body"
			}`,
		},
	}

	for _, tc := range cases {
		userReq := map[string]interface{}{}
		if tc.Username != "" {
			userReq["username"] = tc.Username
		}
		if tc.Email != "" {
			userReq["email"] = tc.Email
		}
		if tc.Password != "" {
			userReq["password"] = tc.Password
		}

		b, _ := json.Marshal(map[string]interface{}{
			"user": userReq,
		})
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/api/users", bytes.NewBuffer(b))

		s.r.ServeHTTP(res, req)

		// then
		s.db.AssertNotCalled(s.T(), "Save", mock.Anything, mock.Anything)
		s.Equal(http.StatusBadRequest, res.Code)
		s.JSONEq(tc.Expected, res.Body.String())
	}
}

func (s *HandlerSuite) TestRegister_FailIfDuplicateError() {
	// given
	username := "zaccoding"
	email := "zaccoding@gmail.com"
	password := "password123"
	matcher := func(acc *model.Account) bool {
		return acc.Username == username && acc.Email == email && MatchesPassword(acc.Password, password) == nil
	}
	s.db.On("Save", mock.Anything, mock.MatchedBy(matcher)).Return(database.ErrKeyConflict)

	// when
	body := map[string]interface{}{
		"user": map[string]interface{}{
			"username": username,
			"email":    email,
			"password": password,
		},
	}
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/api/users", bytes.NewBuffer(b))

	s.r.ServeHTTP(res, req)

	// then
	s.db.AssertCalled(s.T(), "Save", mock.Anything, mock.MatchedBy(matcher))
	s.Equal(http.StatusConflict, res.Code)
	expected := `
	{
	  "code": "DuplicateEntry",
	  "message": "[DuplicateEntry] duplicate email address"
	}`
	s.JSONEq(expected, res.Body.String())
}

func (s *HandlerSuite) TestCurrentUser() {
	// given
	password := "password1"
	encodedPassword, _ := EncodePassword(password)
	acc := model.Account{
		ID:        1,
		Username:  "user1",
		Email:     "user1@gmail.com",
		Password:  encodedPassword,
		Bio:       "user1 bio",
		Image:     "user1 image",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Disabled:  false,
	}
	token := s.getBearerToken(&acc, password)

	// when
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/api/user/me", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	s.r.ServeHTTP(res, req)

	// then
	s.db.AssertCalled(s.T(), "FindByEmail", mock.Anything, acc.Email)
	s.Equal(http.StatusOK, res.Code)
	expected := `
	{
	  "user": {
		"username": "user1",
		"email": "user1@gmail.com",
		"bio": "user1 bio",
		"image": "user1 image"
	  }
	}`
	s.JSONEq(expected, res.Body.String())
}

func (s *HandlerSuite) TestUpdate() {
	// given
	password := "password1"
	encodedPassword, _ := EncodePassword(password)
	acc := model.Account{
		ID:        1,
		Username:  "user1",
		Email:     "user1@gmail.com",
		Password:  encodedPassword,
		Bio:       "user1 bio",
		Image:     "user1 image",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Disabled:  false,
	}
	token := s.getBearerToken(&acc, password)
	s.db.On("Update", mock.Anything, acc.Email, mock.Anything).Return(nil)

	// when
	updateRequest := map[string]interface{}{
		"user": map[string]interface{}{
			"username": "updated-user1",
			"bio":      "updated-bio",
			"image":    "updated-image",
		},
	}
	b, _ := json.Marshal(updateRequest)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/api/user", bytes.NewBuffer(b))
	req.Header.Add("Authorization", "Bearer "+token)

	s.r.ServeHTTP(res, req)

	// then
	s.db.AssertCalled(s.T(), "Update", mock.Anything, acc.Email, mock.MatchedBy(func(a *model.Account) bool {
		return a.Email == acc.Email && a.Username == "updated-user1" && a.Bio == "updated-bio" && a.Image == "updated-image"
	}))
	s.Equal(http.StatusOK, res.Code)
	expected := `
	{
	  "user": {
		"username": "updated-user1",
		"email": "user1@gmail.com",
		"bio": "updated-bio",
		"image": "updated-image"
	  }
	}`
	s.JSONEq(expected, res.Body.String())
}

func (s *HandlerSuite) getBearerToken(acc *model.Account, rawPassword string) string {
	s.db.On("FindByEmail", mock.Anything, acc.Email).Return(acc, nil)
	body := map[string]interface{}{
		"user": map[string]interface{}{
			"email":    acc.Email,
			"password": rawPassword,
		},
	}
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/api/users/login", bytes.NewBuffer(b))
	s.r.ServeHTTP(res, req)

	return gjson.Get(res.Body.String(), "token").String()
}
