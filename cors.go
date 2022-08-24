package zmidd

// 一个gin cors中间件

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultAllowOrigin      = "*"
	DefaultAllowMethods     = "*"
	DefaultAllowHeaders     = "*"
	DefaultAllowCredentials = "*"
	DefaultExposeHeaders    = "true"
)

type Options struct {
	AllowOrigin      string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials string
	ExposeHeaders    string
}

type Option func(s string) *Options

func WithAllowOrigin(allowOrigin string) *Options {
	conf.AllowOrigin = allowOrigin
	return conf
}
func WithAllowMethods(allowMethods string) *Options {
	conf.AllowMethods = allowMethods
	return conf
}
func WithAllowHeaders(allowHeaders string) *Options {
	conf.AllowHeaders = allowHeaders
	return conf
}
func WithAllowCredentials(allowCredentials string) *Options {
	conf.AllowCredentials = allowCredentials
	return conf
}
func WithExposeHeaders(exposeHeaders string) *Options {
	conf.ExposeHeaders = exposeHeaders
	return conf
}

var conf *Options

func init() {
	conf.AllowOrigin = DefaultAllowOrigin
	conf.AllowMethods = DefaultAllowMethods
	conf.AllowHeaders = DefaultAllowHeaders
	conf.AllowCredentials = DefaultAllowCredentials
	conf.ExposeHeaders = DefaultExposeHeaders
}

// Cors
// @Description: 开启跨域控制
// @Param ops: *Options 可选参数 可以设置跨域的配置 默认值为* 允许所有域名访问 允许所有请求方法 跨域允许所有请求头 允许跨域携带cookie
// @Return gin.HandlerFunc 返回一个gin.HandlerFunc
// 可以直接使用gin.Use(Cors()) 或者在路由中使用Cors() 即可
// 如果需要设置跨域的配置 可以使用 WithAllowOrigin WithAllowMethods WithAllowHeaders WithAllowCredentials WithExposeHeaders
func Cors(ops ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", conf.AllowOrigin)
			c.Header("Access-Control-Allow-Methods", conf.AllowMethods)
			c.Header("Access-Control-Allow-Headers", conf.AllowHeaders)
			c.Header("Access-Control-Allow-Credentials", conf.AllowCredentials)
			c.Header("Access-Control-Expose-Headers", conf.ExposeHeaders)
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
