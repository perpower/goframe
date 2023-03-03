// 统一错误处理
// Author: syswen
package errors

import "fmt"

// 定义错误返回结构体
type OutError struct {
	Code int         `json:"code" xml:"code" yml:"code"` // 错误码
	Msg  string      `json:"msg" xml:"msg" yml:"msg"`    // 错误信息
	Data interface{} `json:"data" xml:"data" yml:"data"` // 返回数据
}

var emptyStruct = struct{}{}

var (
	SUCCESS_CODE = OutError{0, "success", emptyStruct}
	ERROR_CODE   = OutError{-1, "failed", emptyStruct}
	ERROR_1001   = OutError{1001, "签名失效", emptyStruct}
	ERROR_1002   = OutError{1002, "签名错误", emptyStruct}
	ERROR_2001   = OutError{2001, "数据记录不存在", emptyStruct}
	ERROR_2002   = OutError{2002, "数据重复", emptyStruct}
	ERROR_2003   = OutError{2003, "创建/更新数据失败", emptyStruct}
	ERROR_3000   = OutError{3000, "参数验证不通过", emptyStruct}
	ERROR_3001   = OutError{3001, "操作过于频繁请稍后再试", emptyStruct}
	ERROR_3002   = OutError{3002, "无操作权限", emptyStruct}
	ERROR_3003   = OutError{3003, "上传失败", emptyStruct}
	ERROR_3004   = OutError{3004, "数据格式不正确", emptyStruct}
	ERROR_3005   = OutError{3005, "提交的数据不符合字典约束范围值", emptyStruct}
	ERROR_3006   = OutError{3006, "提交的数据校验不通过，验证失败", emptyStruct}
	ERROR_3054   = OutError{3054, "系统繁忙,请稍后再试", emptyStruct}
	ERROR_4001   = OutError{4001, "未授权", emptyStruct}
	ERROR_4002   = OutError{4002, "未知错误", emptyStruct}
	ERROR_4004   = OutError{4002, "页面未定义", emptyStruct}
	ERROR_5000   = OutError{5000, "服务器异常", emptyStruct}
	ERROR_9000   = OutError{9000, "账户授权Token值已过期请重新获取", emptyStruct}
	ERROR_9001   = OutError{9001, "账户被禁用请联系管理员", emptyStruct}
	ERROR_9002   = OutError{9002, "账号/密码错误请检查后重试", emptyStruct}
	ERROR_9003   = OutError{9003, "账号已存在", emptyStruct}
	ERROR_9004   = OutError{9004, "账号/密码错误请检查后重试", emptyStruct}
	ERROR_9005   = OutError{9005, "密码错误次数过多请稍后再试", emptyStruct}
	ERROR_9006   = OutError{9006, "账户信息异常", emptyStruct}
)

func (e *OutError) Error() string {
	return e.Msg
}

// New creates and returns an error code.
// Note that it returns an interface object of Code.
// code: int
// msg: string
// data: interface{}
func New(code int, msg string, data interface{}) *OutError {
	if data == nil {
		data = emptyStruct // nil 转空结构体
	}
	return &OutError{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Newf creates and returns an error code with format.
// Note that it returns an interface object of Code.
// code: int
// msg: string
// data: interface{}
// msgArgs: slice interface{}
func Newf(code int, msg string, data interface{}, msgArgs ...interface{}) *OutError {
	if data == nil {
		data = emptyStruct // nil 转空结构体
	}
	return &OutError{
		Code: code,
		Msg:  fmt.Sprintf(msg, msgArgs...),
		Data: data,
	}
}

// GetCode returns the integer number of current error code.
func (c *OutError) GetCode() int {
	return c.Code
}

// GetMsg returns the brief message for current error code.
func (c *OutError) GetMsg() string {
	return c.Msg
}

// GetData returns the detailed information of current error code,
// which is mainly designed as an extension field for error code.
func (c *OutError) GetData() interface{} {
	return c.Data
}

// GetString returns current error code as a string.
func (c *OutError) GetString() string {
	if c.Data != nil {
		return fmt.Sprintf(`%d:%s %v`, c.Code, c.Msg, c.Data)
	}
	if c.Msg != "" {
		return fmt.Sprintf(`%d:%s`, c.Code, c.Msg)
	}
	return fmt.Sprintf(`%d`, c.Code)
}
