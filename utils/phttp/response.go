package phttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perpower/goframe/utils/errors"
)

// OutJson 定义统一输出json格式内容的方法
// obj：interface{} 待输出的内容结构体
func OutJson(c *gin.Context, obj interface{}) (typ interface{}) {
	switch objType := obj.(type) {
	case errors.OutError:
		c.JSON(http.StatusOK, obj)
		typ = objType
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": errors.SUCCESS_CODE.Code,
			"msg":  errors.SUCCESS_CODE.Msg,
			"data": obj,
		})
		typ = objType
	}
	c.Abort()
	return typ
}
