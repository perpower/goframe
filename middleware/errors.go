// errors错误处理中间件
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perpower/goframe/utils/errors"
)

// ErrorHandle
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先调用c.Next()执行后面的中间件
		// 所有中间件及router处理完毕后从这里开始执行
		// 检查c.Errors中是否有错误

		for _, e := range c.Errors {
			err := e.Err
			// 若是自定义的错误则将code、msg, data返回
			if errInfo, ok := err.(*errors.OutError); ok {
				c.JSON(http.StatusOK, gin.H{
					"code": errInfo.Code,
					"msg":  errInfo.Msg,
					"data": errInfo.Data,
				})
			} else {
				// 若非自定义错误则返回详细错误信息err.Error()
				c.JSON(http.StatusOK, gin.H{
					"code": errors.ERROR_5000.Code,
					"msg":  "服务器异常",
					"data": err.Error(),
				})
			}
			return // 只要检查到一个错误就返回
		}
	}
}
