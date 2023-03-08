// 日志记录中间件
package middleware

import (
	"go-framework/global"

	"github.com/perpower/goframe/utils/logger"

	"github.com/gin-gonic/gin"
)

// LoggerHandle
// conf: interface{}
func LoggerHandle(conf interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		plog := logger.InitLogger(c, global.LogPlatform, conf)
		global.Plog = plog
		c.Next()
	}
}
