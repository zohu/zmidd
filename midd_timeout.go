package zmidd

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zohu/zutils/zbuffpool"
	"net/http"
	"strings"
	"time"
)

type SimpleBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w SimpleBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// Timeout
// @Description: 超时中间件
// @param timeout
// @param whitelist
// @return gin.HandlerFunc
func Timeout(timeout time.Duration, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果超时时间为0，则不启用超时中间件
		if timeout == time.Duration(0) {
			c.Next()
			return
		}
		// 如果请求路径在白名单中，则不启用超时中间件
		for _, prefix := range whitelist {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Next()
				return
			}
		}
		// 创建一个缓冲区
		buffer := zbuffpool.GetBuff()
		// 创建一个新的响应写入器
		blw := &SimpleBodyWriter{body: buffer, ResponseWriter: c.Writer}
		// 替换原有的响应写入器
		c.Writer = blw
		// 创建一个新的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		// 替换原有的上下文
		c.Request = c.Request.WithContext(ctx)
		// 创建一个完成通道
		finish := make(chan struct{})
		// 创建一个子协程
		go func() {
			// 执行后续的中间件
			c.Next()
			// 向完成通道发送完成信号
			finish <- struct{}{}
		}()
		// 等待完成通道或者超时通道
		select {
		case <-ctx.Done():
			// 如果超时发生，则向客户端返回超时状态码
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			// 中断请求
			c.Abort()
			// 超时发生, 通知子协程退出
			cancel()
			// 如果超时的话, buffer无法主动清除, 只能等待GC回收
		case <-finish:
			// 结果只会在主协程中被写入
			_, _ = blw.ResponseWriter.Write(buffer.Bytes())
			// 释放缓冲区
			zbuffpool.PutBuff(buffer)
			// 释放上下文
			cancel()
		}
	}
}
