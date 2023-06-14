package redis

import (
	"reflect"
)

type Config struct {
	Address     string // 地址 ip:port
	Username    string // redis6.0版本以上开始提供Redis ACL,用户名+密码一起使用
	Password    string // 密码
	Database    int    // 数据库
	UseTLS      bool   // 是否启用tls
	MaxIdle     int    // 允许闲置的最大连接数(0表示不限制)
	MaxActive   int    // 最大连接数量限制(0表示不限制)
	IdleTimeout int    // 连接最大空闲时间
}

var (
	defaultDatabase    = 0     // 默认数据库
	defaultUseTLS      = false // 默认是否使用TLS
	defaultMaxIdle     = 10
	defaultMaxActive   = 100
	defaultIdleTimeout = 10
)

func SetConfig(conf Config) Config {
	cnf := reflect.TypeOf(conf)

	if _, ok := cnf.FieldByName("Database"); !ok {
		conf.Database = defaultDatabase
	}

	if _, ok := cnf.FieldByName("UseTLS"); !ok {
		conf.UseTLS = defaultUseTLS
	}

	if _, ok := cnf.FieldByName("MaxIdle"); !ok {
		conf.MaxIdle = defaultMaxIdle
	}

	if _, ok := cnf.FieldByName("MaxActive"); !ok {
		conf.MaxActive = defaultMaxActive
	}

	if _, ok := cnf.FieldByName("IdleTimeout"); !ok {
		conf.IdleTimeout = defaultIdleTimeout
	}

	// c.config = &conf

	return conf
}
