package middleware

import (
	"gin-rest-api-example/pkg/trace"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestIDMiddleware(t *testing.T) {
	type assertFunc func(id string)
	cases := []struct {
		Name       string
		RequestId  string
		AssertFunc assertFunc
	}{
		{
			Name: "empty request id",
			AssertFunc: func(id string) {
				assert.Equal(t, len(id), 36)
			},
		}, {
			Name:      "use requested value",
			RequestId: "custom",
			AssertFunc: func(id string) {
				assert.Equal(t, "custom", id)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			s := setupRouterWithHandler(func(c *gin.Engine) {
				c.Use(RequestIDMiddleware())
			}, func(c *gin.Context) {
				requestId := trace.RequestIDFromContext(c)
				tc.AssertFunc(requestId)
			})

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "http://localhost/foo", nil)
			if tc.RequestId != "" {
				req.Header.Add(XRequestIdKey, tc.RequestId)
			}

			// when then
			s.ServeHTTP(res, req)
		})
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	now := time.Now()
	s := setupRouterWithHandler(func(c *gin.Engine) {
		c.Use(TimeoutMiddleware(time.Second))
	}, func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		assert.True(t, ok)
		assert.LessOrEqual(t, deadline.Sub(now).Milliseconds(), time.Second.Milliseconds())
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost/foo", nil)

	// when then
	s.ServeHTTP(res, req)
}

func setupRouterWithHandler(middlewareFunc func(c *gin.Engine), handler func(c *gin.Context)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	middlewareFunc(r)
	r.GET("/foo", handler)
	return r
}
