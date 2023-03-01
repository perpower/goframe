// zap日志记录组件
package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/perpower/goframe/funcs"
	"github.com/perpower/goframe/funcs/ptime"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var ctx *gin.Context

// 日志内容自定义数据key=>value结构体
type ExtendFields struct {
	Key   string
	Value interface{}
}

// request请求信息结构体
type requestInfo struct {
	RequestTime string //请求时间
	RequestURL  string //请求地址
	RequestUA   string //UserAgent
	RequestIP   string //请求IP
	RequestBody string //请求body
}

type LogFileConfig struct {
	FilePath   string //文件路径,不包含文件名
	MaxSize    int    //单个文件大小,单位M
	MaxBackups int    //最大保留旧日志文件数量
	MaxAge     int    //日志文件最长保留天数
	Compress   bool   //是否使用gzip压缩已旋转的日志文件,默认是不执行压缩
}

func InitLogger(c *gin.Context, conf LogFileConfig) {
	encoder := getJsonEncoder()

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//info文件WriteSyncer
	infoFileWriteSyncer := getInfoWriterSyncer(conf)
	//error文件WriteSyncer
	errorFileWriteSyncer := getErrorWriterSyncer(conf)

	//生成core
	//同时输出到控制台 和 指定的日志文件中
	infoFileCore := zapcore.NewCore(getConsoleEncoder(), infoFileWriteSyncer, lowPriority)
	errorFileCore := zapcore.NewCore(encoder, errorFileWriteSyncer, highPriority)

	//将infocore 和 errcore 加入core切片
	var coreArr []zapcore.Core
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	//生成Logger
	Logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller() 显示文件名 和 行号
	ctx = c
}

func newEncoderConfig(levelEncoder zapcore.LevelEncoder) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caler",
		MessageKey:     "Msg",
		StacktraceKey:  "Stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    levelEncoder, //zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,  //指定时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   shortCallerEncoder, //zapcore.ShortCallerEncoder,
	}
}

// 日志时间格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(ptime.Format_date_time))
}

func shortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", caller.TrimmedPath()))
}

// Encoder 获取Json编码器
func getJsonEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(newEncoderConfig(zapcore.CapitalLevelEncoder))
}

// Encoder 获取Console编码器
func getConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(newEncoderConfig(zapcore.CapitalColorLevelEncoder))
}

// Info级别日志输出路径,同时在控制台，日志文件中输出
func getInfoWriterSyncer(conf LogFileConfig) zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能
	infoLumberIO := &lumberjack.Logger{
		Filename:   conf.FilePath + "/info.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,   // days
		Compress:   conf.Compress, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(infoLumberIO), zapcore.AddSync(os.Stdout))
}

// Error级别日志输出路径,同时在控制台，日志文件中输出
func getErrorWriterSyncer(conf LogFileConfig) zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   conf.FilePath + "/error.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,   // days
		Compress:   conf.Compress, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberWriteSyncer), zapcore.AddSync(os.Stdout))
}

// 二次封装日志级别记录方法
func Debug(msg string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Debug(msg, fields...)
}

func Info(msg string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Info(msg, fields...)
}

func Warn(msg string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Warn(msg, fields...)
}

func Error(msg string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Error(msg, fields...)
}

func Panic(format string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Panic(format, fields...)
}

func Fatal(format string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	Logger.Fatal(format, fields...)
}

// 处理额外数据
func convertFields(filedSlice ...ExtendFields) zapcore.Field {
	fileds := []string{}
	for _, extend := range filedSlice {
		fileds = append(fileds, fmt.Sprintf("%+v", extend))
	}

	return zap.Strings("ExtraDatas", fileds)
}

// 默认补充Request 请求基础数据
func convertRequestInfo() zapcore.Field {
	requestInfo := requestInfo{
		RequestTime: ptime.TimestampStr(),
		RequestURL:  ctx.Request.Method + "  " + ctx.Request.Host + ctx.Request.RequestURI,
		RequestUA:   ctx.Request.UserAgent(),
		RequestIP:   ctx.ClientIP(),
	}

	requestBody, _ := io.ReadAll(ctx.Request.Body)
	if requestBody != nil {
		requestInfo.RequestBody = funcs.Bytes2String(requestBody)
	}
	// 通过 ioutil.ReadAll() 来读取完 body 内容后，body 就为空了，把字节流重新放回 body 中
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	return zap.String("RequestInfo", fmt.Sprintf("%+v", requestInfo))
}
