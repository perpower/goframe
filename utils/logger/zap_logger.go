// zap日志记录组件
package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/perpower/goframe/funcs/ptime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFileConfig struct {
	RootDir    string // 日志文件存放跟目录
	FilePath   string //文件路径,不包含文件名
	MaxSize    int    //单个文件大小,单位M
	MaxBackups int    //最大保留旧日志文件数量
	MaxAge     int    //日志文件最长保留天数
	Compress   bool   //是否使用gzip压缩已旋转的日志文件,默认是不执行压缩
}

func InitLocal(conf LogFileConfig) {
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
		Filename:   conf.RootDir + conf.FilePath + "/access.log",
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
		Filename:   conf.RootDir + conf.FilePath + "/error.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,   // days
		Compress:   conf.Compress, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberWriteSyncer), zapcore.AddSync(os.Stdout))
}

// CreateFileLog 创建日志文件
// level: string 错误等级
// msg: string 消息文本
// filedSlice: []ExtendFields  额外参数
func CreateFileLog(level, msg string, filedSlice ...ExtendFields) {
	fields := []zapcore.Field{convertRequestInfo()}
	if len(filedSlice) > 0 {
		fields = append(fields, convertFields(filedSlice...))
	}
	switch level {
	case "debug":
		Logger.Debug(msg, fields...)
	case "info":
		Logger.Info(msg, fields...)
	case "warn":
		Logger.Warn(msg, fields...)
	case "error":
		Logger.Error(msg, fields...)
	case "panic":
		Logger.Panic(msg, fields...)
	case "fatal":
		Logger.Fatal(msg, fields...)
	}
}

// convertFields 处理额外数据
func convertFields(filedSlice ...ExtendFields) zapcore.Field {
	fileds := []string{}
	for _, extend := range filedSlice {
		fileds = append(fileds, fmt.Sprintf("%+v", extend))
	}

	return zap.Strings("ExtraDatas", fileds)
}

// convertRequestInfo 默认补充Request 请求基础数据
func convertRequestInfo() zapcore.Field {
	requestInfo := requestInfo()
	return zap.String("RequestInfo", fmt.Sprintf("%+v", requestInfo))
}
