// panic错误捕获处理
package middleware

import (
	"github.com/perpower/goframe/funcs"
	"github.com/perpower/goframe/utils/alarm"

	"github.com/gin-gonic/gin"
)

func ErrorHandle(appName, emailTpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//邮件告警通知
				alarm.Email(c, appName, funcs.GetRootPath()+"/"+emailTpl, err)
			}
		}()

		c.Next()
	}
}
