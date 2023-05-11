package plog

import (
	"bytes"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/perpower/goframe/funcs/normal"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/utils/pelastic"
	"go.uber.org/zap"
)

// 日志内容自定义数据key=>value结构体
type ExtendFields struct {
	Key   string
	Value interface{}
}

// request请求信息结构体
type requestFields struct {
	RequestTime   string              `json:"requestTime"`   // 请求时间
	RequestMethod string              `json:"requestMethod"` // 请求方式
	RequestProto  string              `json:"requestProto"`  // 请求协议
	RequestHost   string              `json:"requestHost"`   // 主机地址
	RequestUri    string              `json:"requestUri"`    // 请求地址
	UserAgent     string              `json:"userAgent"`     // UserAgent
	ClientIp      string              `json:"clientIp"`      // 请求IP
	Headers       map[string][]string `json:"headers"`       // 请求header
	Refer         string              `json:"refer"`         // 请求 refer
	RequestBody   string              `json:"requestBody"`   // 请求body
}

// 定义日志输出的接口方法
type StandLog interface {
	Debug(cate, msg string, filedSlice ...ExtendFields) // debug 级别日志
	Info(cate, msg string, filedSlice ...ExtendFields)  // info 级别日志
	Warn(cate, msg string, filedSlice ...ExtendFields)  // warn 级别日志
	Error(cate, msg string, filedSlice ...ExtendFields) // error 级别日志
	Panic(cate, msg string, filedSlice ...ExtendFields) // panic 级别日志
	Fatal(cate, msg string, filedSlice ...ExtendFields) // fatal 级别日志
}

type Output struct {
	platform string //日志存储平台
}

var (
	Logger *zap.Logger
	ctx    *gin.Context
)

// InitLogger 日志服务初始化
// conf: interface{}
func InitLogger(c *gin.Context, platform string, conf interface{}) *Output {
	if _, ok := conf.(LogFileConfig); ok {
		InitLocal(conf.(LogFileConfig)) //初始化日志组件
		defer Logger.Sync()
	} else if _, ok := conf.(pelastic.ElastiConfig); ok {
		InitElastic(conf.(pelastic.ElastiConfig))
	}

	out := &Output{
		platform: platform,
	}
	ctx = c
	return out
}

// Request 请求基础数据
func requestInfo() requestFields {
	requestInfo := requestFields{
		RequestTime:   ptime.TimestampMilliStr(),
		RequestMethod: ctx.Request.Method,
		RequestProto:  ctx.Request.Proto,
		RequestHost:   ctx.Request.Host,
		RequestUri:    ctx.Request.RequestURI,
		UserAgent:     ctx.Request.UserAgent(),
		ClientIp:      ctx.ClientIP(),
		Headers:       ctx.Request.Header,
		Refer:         ctx.Request.Referer(),
	}

	requestBody, _ := io.ReadAll(ctx.Request.Body)
	if requestBody != nil {
		requestInfo.RequestBody = normal.Bytes2String(requestBody)
	}
	// 通过 ioutil.ReadAll() 来读取完 body 内容后，body 就为空了，把字节流重新放回 body 中
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	return requestInfo
}

func (c *Output) Debug(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("debug", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("debug", cate, msg, filedSlice...)
	}
}

func (c *Output) Info(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("info", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("info", cate, msg, filedSlice...)
	}
}

func (c *Output) Warn(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("warn", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("warn", cate, msg, filedSlice...)
	}
}

func (c *Output) Error(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("error", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("error", cate, msg, filedSlice...)
	}
}

func (c *Output) Panic(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("panic", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("panic", cate, msg, filedSlice...)
	}
}

func (c *Output) Fatal(cate, msg string, filedSlice ...ExtendFields) {
	switch c.platform {
	case "file":
		CreateFileLog("fatal", msg, filedSlice...)
	case "elasticSearch":
		CreateElasticLog("fatal", cate, msg, filedSlice...)
		os.Exit(1)
	}
}
