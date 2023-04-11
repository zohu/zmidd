package zmidd

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/zohu/zlog"
	"github.com/zohu/zutils"
	"github.com/zohu/zutils/zbuffpool"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoggerConfig struct {
	LogTag      string `yaml:"log_tag" mapstructure:"log_tag"`
	MaxBody     int    `yaml:"max_body" mapstructure:"max_body"`
	MaxResponse int    `yaml:"max_response" mapstructure:"max_response"`
	MaxFile     uint   `yaml:"max_file" mapstructure:"max_file"`
	Path        string `yaml:"path" mapstructure:"path"`
	PrintThird  bool   `yaml:"print_third" mapstructure:"print_third"`
}

type LogLinkField struct {
	Rid   string `json:"rid"`
	Time  int64  `json:"time"`
	Delay int64  `json:"delay"`
	Tag   string `json:"tag"`
	Auth  string `json:"auth"`
	Head  string `json:"head"`
	Param string `json:"param"`
	Data  string `json:"data"`
	Error string `json:"error"`
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(conf LoggerConfig, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 健康检查不记录日志
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}
		// 如果请求路径在白名单中，则不记录日志
		for _, prefix := range whitelist {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Next()
				return
			}
		}
		// 创建一个缓冲区
		buffer := zbuffpool.GetBuff()
		// 创建一个新的响应写入器
		blw := &CustomResponseWriter{body: buffer, ResponseWriter: c.Writer}
		// 替换原有的响应写入器
		c.Writer = blw

		// 构建日志
		lf := new(LogLinkField)
		lf.Tag = conf.LogTag
		lf.Rid = zutils.FirstValue(c.Request.Header.Get("X-Request-ID"), GetRequestId(c))
		lf.Time = time.Now().UnixMicro()
		ua := user_agent.New(c.Request.UserAgent())
		brow, browVersion := ua.Browser()
		lf.Head = fmt.Sprintf(
			"%s %s %s ip=%s %s-%s",
			c.Request.Method,   // 请求方法
			c.Request.URL.Path, // 请求路径
			lf.Rid,             // 请求ID
			c.ClientIP(),       // ip
			brow,               // 浏览器
			browVersion,        // 浏览器版本
		)
		// 客户端版本
		vsn := c.GetHeader("X-Client-Version")
		if vsn != "" {
			lf.Head = fmt.Sprintf("%s vsn=%s", lf.Head, vsn)
		}
		// 请求参数
		if c.ContentType() == gin.MIMEMultipartPOSTForm {
			lf.Param = fmt.Sprintf("Form query: %v", c.Request.URL.RawQuery)
		} else {
			data, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			lf.Param = zutils.FirstValue(string(data), c.Request.URL.RawQuery)
			if len([]rune(lf.Param)) > conf.MaxBody {
				lf.Param = string([]rune(lf.Param)[:conf.MaxBody]) + "..."
			}
		}

		// 执行请求
		start := time.Now()
		c.Next()
		// 获取返回状态和返回值
		lf.Head = fmt.Sprintf("%d %s", c.Writer.Status(), lf.Head)
		// 获取返回值
		lf.Data = blw.body.String()
		if len([]rune(lf.Data)) > conf.MaxResponse {
			lf.Data = string([]rune(lf.Data)[:conf.MaxResponse]) + "..."
		}
		// 计算耗时
		lf.Delay = time.Since(start).Milliseconds()

		// 记录日志
		logs := fmt.Sprintf("%s %dms %s", lf.Head, lf.Delay, lf.Tag)
		if lf.Param != "" {
			logs = fmt.Sprintf("%s param=%s", logs, lf.Param)
		}
		if lf.Data != "" {
			logs = fmt.Sprintf("%s data=%s", logs, lf.Data)
		}
		if len(c.Errors) > 0 {
			lf.Error = strings.Replace(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n", "", -1)
			zlog.WithOptions(zap.WithCaller(false)).Sugar().Warnf("%s error=%s", logs, lf.Error)
		} else {
			if c.Writer.Status() >= http.StatusBadRequest {
				lf.Error = blw.body.String()
				zlog.WithOptions(zap.WithCaller(false)).Sugar().Warnf("%s error=%s", logs, lf.Error)
			} else {
				auth, ok := c.Get("auth")
				if ok {
					lf.Auth, _ = zutils.StructToString(auth)
					logs = fmt.Sprintf("%s auth=%s", logs, lf.Auth)
				}
				zlog.WithOptions(zap.WithCaller(false)).Sugar().Infof(logs)
			}
		}
	}
}
