// 程序错误告警--可扩展邮件，短信，微信等等告警方式
package alarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-framework/configs"
	"io"
	"runtime/debug"
	"strings"

	"github.com/perpower/goframe/funcs/normal"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/utils/mailer"

	"github.com/gin-gonic/gin"
)

type errorString struct {
	s string
}

// 告警内容结构体
type errorInfo struct {
	AppName     string // 系统名称
	ErrorMsg    string // 错误信息
	RequestTime string // 请求时间
	RequestURL  string // 请求地址
	RequestBody string
	UserAgent   string   // UserAgent
	ClientIp    string   // 请求IP
	Headers     string   // 请求header
	Refer       string   // 请求 refer
	DebugStack  []string // 错误跟踪信息
}

func (e *errorString) Error() string {
	return e.s
}
func Info(c *gin.Context, appName string, err interface{}) error {
	alarm(c, appName, "info", err)
	return &errorString{fmt.Sprintf("%v", err)}
}

// Email 发邮件
// appName: string 系统名称
// emailTpl: string 邮件模板文件路径
// err: interface{} 错误信息
func Email(c *gin.Context, appName, emailTpl string, err interface{}) error {
	ErrorMsg := alarm(c, appName, "email", err)

	subject := fmt.Sprintf("【错误告警】- %s 项目出错了！", appName)
	body, err := mailer.GetTplContentByFile(emailTpl, ErrorMsg)
	if err == nil {
		mailer.Send(mailer.EmailSererConfig{
			ServerAddress: configs.EmailAccount["serverAddress"].(string),
			Port:          configs.EmailAccount["port"].(int),
			Username:      configs.EmailAccount["username"].(string),
			Password:      configs.EmailAccount["password"].(string),
		}, mailer.EmailConfig{
			To:      []string{configs.EmailAccount["username"].(string)},
			Subject: subject,
			Body:    body,
		})
	}
	return &errorString{fmt.Sprintf("%v", err)}
}

// alarm 告警方法
// appName: string 系统名称
// level: string 错误等级
// err: interface{} 错误信息
func alarm(c *gin.Context, appName, level string, err interface{}) (ErrorMsg errorInfo) {
	DebugStack := strings.Split(string(debug.Stack()), "\n")

	headers, _ := json.Marshal(c.Request.Header)
	ErrorMsg = errorInfo{
		AppName:     appName,
		ErrorMsg:    fmt.Sprintf("%s", err),
		RequestTime: ptime.TimestampStr(),
		RequestURL:  c.Request.Method + "  " + c.Request.Host + c.Request.RequestURI,
		UserAgent:   c.Request.UserAgent(),
		ClientIp:    c.ClientIP(),
		Headers:     normal.Bytes2String(headers),
		Refer:       c.Request.Referer(),
		DebugStack:  DebugStack,
	}

	requestBody, _ := io.ReadAll(c.Request.Body)
	if requestBody != nil {
		ErrorMsg.RequestBody = normal.Bytes2String(requestBody)
	}
	// 通过 ioutil.ReadAll() 来读取完 body 内容后，body 就为空了，把字节流重新放回 body 中
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	json.Marshal(ErrorMsg)

	return ErrorMsg
}
