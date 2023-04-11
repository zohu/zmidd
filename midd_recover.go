package zmidd

import (
	"github.com/gin-gonic/gin"
	"github.com/zohu/zflag"
	"github.com/zohu/zlog"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					zlog.Errorf("系统错误 %s", e.Error())
					zflag.Done(c, zflag.ResponseErr(zflag.ErrNil.WithMessage(e.Error())))
					return
				}
				zlog.Errorf("系统错误:%v", err)
				zflag.Done(c, zflag.ResponseErr(zflag.ErrNil))
			}
		}()
		c.Next()
	}
}
