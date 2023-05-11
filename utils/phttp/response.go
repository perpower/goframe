package phttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perpower/goframe/utils/perrors"
)

// OutJson 定义统一响应json格式内容的方法
// obj：interface{} 待响应的内容结构体
func OutJson(c *gin.Context, obj interface{}) {
	if _, ok := obj.(perrors.OutError); ok {
		c.JSON(http.StatusOK, obj)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": perrors.SUCCESS_CODE.Code,
		"msg":  perrors.SUCCESS_CODE.Msg,
		"data": obj,
	})
}

// ThrowError
// err: perrors.OutError 标准错误结构体
func ThrowError(c *gin.Context, err perrors.OutError) {
	c.Error(err)
}

// ThrowErrorMsg
// msg: string 错误信息
// args: []interface{} 格式化错误信息对应的参数
func ThrowErrorMsg(c *gin.Context, msg string, args ...interface{}) {
	var err perrors.OutError
	if len(args) > 0 {
		err = perrors.Newf(perrors.ERROR_CODE.Code, msg, nil, args...)
	} else {
		err = perrors.New(perrors.ERROR_CODE.Code, msg, nil)
	}

	c.Error(err)
}
