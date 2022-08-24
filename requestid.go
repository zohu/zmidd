package zmidd

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func RequestId() gin.HandlerFunc {
	return requestid.New()
}

func GetRequestId(c *gin.Context) string {
	return requestid.Get(c)
}
