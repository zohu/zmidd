package zmidd

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/zohu/zutils"
)

const headerXRequestID = "X-Request-ID"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = xid.New().String()
			c.Request.Header.Add(headerXRequestID, rid)
		}
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

func GetRequestId(c *gin.Context) string {
	return zutils.FirstValue(
		c.Writer.Header().Get(headerXRequestID),
		c.Request.Header.Get(headerXRequestID),
		c.GetHeader(headerXRequestID),
	)
}
