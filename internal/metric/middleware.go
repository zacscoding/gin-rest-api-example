package metric

import (
	"gin-rest-api-example/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

func Route(r *gin.Engine) {
	r.GET("metric", gin.WrapH(promhttp.Handler()))
}

func MetricsMiddleware(mp *MetricsProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Now().Sub(start)
		logging.DefaultLogger().Infof("## Called code:%d, method:%s, path:%s, elapsed:%v", c.Writer.Status(), c.Request.Method, c.FullPath(), elapsed)
		var (
			code   = c.Writer.Status()
			method = c.Request.Method
			path   = c.FullPath()
		)
		mp.RecordApiCount(code, method, path)
		mp.RecordApiLatency(code, method, path, elapsed)
	}
}
