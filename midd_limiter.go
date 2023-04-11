package zmidd

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/zohu/zlog"
	"time"

	limits "github.com/gin-contrib/size"
)

// RateLimiter
// @Description: 请求速率限制
// @return gin.HandlerFunc
func RateLimiter(rate float64) gin.HandlerFunc {
	if rate > 0 {
		zlog.Infof("接口限速 %f", rate)
		lmt := tollbooth.NewLimiter(rate, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
		lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
		return func(c *gin.Context) {
			httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
			if httpError != nil {
				c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
				c.Abort()
			} else {
				c.Next()
			}
		}
	} else {
		return func(c *gin.Context) {
			c.Next()
		}
	}
}

// SizeLimiter
// @Description: 请求体大小限制
// @return gin.HandlerFunc
func SizeLimiter(size int64) gin.HandlerFunc {
	if size > 0 {
		zlog.Infof("请求体限制 %dM", size)
		return limits.RequestSizeLimiter(size * 1024 * 1024)
	} else {
		return func(c *gin.Context) {
			c.Next()
		}
	}
}
