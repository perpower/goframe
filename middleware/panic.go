// panic错误捕获处理
package middleware

import (
	"github.com/perpower/goframe/utils/palarm"
	"github.com/perpower/goframe/utils/pmailer"

	"github.com/gin-gonic/gin"
)

// PanicHandle 捕获服务允许过程中发生的panic错误
// appName: string 系统服务名
// emailServerConfig: mailer.EmailSererConfig 发件服务配置
// receivers: []string 告警邮件收件人地址
// emailTpl: string 邮件模板
func PanicHandle(appName string, emailServerConfig pmailer.EmailSererConfig, receivers []string, emailTpl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//邮件告警通知
				palarm.Email(c, appName, emailServerConfig, receivers, "./"+emailTpl, err)
			}
		}()

		c.Next()
	}
}
