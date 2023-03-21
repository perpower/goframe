// 日志记录中间件
package middleware

import (
	"github.com/perpower/goframe/utils/logger"

	"github.com/gin-gonic/gin"
)

// LoggerHandle
// plog: *logger.Output 日志服务指针
// logPlatform: string 平台名
// conf: interface{} 日志平台配置
func LoggerHandle(plog *logger.Output, logPlatform string, conf interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		*plog = *logger.InitLogger(c, logPlatform, conf)
		c.Next()
	}
}
