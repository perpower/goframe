// 日志记录中间件
package middleware

import (
	"github.com/perpower/goframe/utils/logger"

	"github.com/gin-gonic/gin"
)

func LoggerHandle(conf logger.LogFileConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InitLogger(c, conf) //初始化日志组件
		defer logger.Logger.Sync()

		c.Next()
	}
}
