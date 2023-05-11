// 日志记录中间件
package middleware

import (
	"github.com/perpower/goframe/utils/plog"

	"github.com/gin-gonic/gin"
)

// LoggerHandle
// plog: *plog.Output 日志服务指针
// logPlatform: string 平台名
// conf: interface{} 日志平台配置
func LoggerHandle(ploger *plog.Output, logPlatform string, conf interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		*ploger = *plog.InitLogger(c, logPlatform, conf)
		c.Next()
	}
}
