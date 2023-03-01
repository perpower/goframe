package ptime

import (
	"time"
)

type Time struct {
	wrapper
}

var (
	Format_date_time = "2006-01-02 15:04:05"
	Format_date      = "2006-01-02"
	Format_iso       = "2006-01-02T15:04:05-07:00"
	Format_rfc       = "Mon, 02 Jan 06 15:04 MST"
)

// 返回当前时间的秒级时间戳
func Timestamp() int64 {
	return Now().Unix()
}

// 返回当前时间的毫秒级时间戳
func TimestampMilli() int64 {
	return Now().UnixMilli()
}

// 返回当前时间的微秒级时间戳
func TimestampMicro() int64 {
	return Now().UnixMicro()
}

// 返回当前时间的纳秒级时间戳
func TimestampNano() int64 {
	return Now().UnixNano()
}

// Date returns current date in string like "2006-01-02".
func Date() string {
	return time.Now().Format(Format_date)
}

// Date returns current datetime in string like "2006-01-02 15:04:05".
func DateTime() string {
	return time.Now().Format(Format_date_time)
}

// ISO8601 returns current datetime in ISO8601 format like "2006-01-02T15:04:05-07:00".
func ISO8601() string {
	return time.Now().Format(Format_iso)
}

// RFC822 returns current datetime in RFC822 format like "Mon, 02 Jan 06 15:04 MST".
func RFC822() string {
	return time.Now().Format(Format_rfc)
}

// 毫秒时间戳间戳转换成日期格式Y-m-d H:i:s
func UnixToDate(timestamp int64) string {
	tp := time.UnixMilli(timestamp)
	return tp.Format(Format_date_time)
}

// 日期转换成毫秒时间戳 2006-01-02 15:04:05
func DateToUnix(str string) int64 {
	template := Format_date_time
	tp, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return tp.UnixMilli()
}
