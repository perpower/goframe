// panic错误捕获处理
package middleware

import (
	"github.com/perpower/goframe/utils/alarm"

	"github.com/gin-gonic/gin"
)

// PanicHandle 捕获服务允许过程中发生的panic错误
func PanicHandle(appName, emailTpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//邮件告警通知
				alarm.Email(c, appName, "./"+emailTpl, err)
			}
		}()

		c.Next()
	}
}
