// redis功能封装--单例模式
// 特别说明：该功能包要求redis版本使用7.0.0以上，引入了redis7.0的一些新特性，否则有些方法会不支持
package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	config      *Config
	conn        redis.Conn
	Db          *Rdb
	Scan        *Rscan
	String      *Rstring
	Expire      *Rexpire
	Hash        *Rhash
	Set         *Rset
	Zset        *Rzset
	List        *Rlist
	Geo         *Rgeo
	Transaction *Rtransaction
	Hyper       *Rhyper
	Bit         *Rbit
}

var (
	defaultFlushdbMode = "SYNC" // FLUSHDB 默认清空模式
	defaultCursor      = 0      // SCAN默认起始游标
	defaultScanNum     = 10     // 单次scan迭代数量
)

// 返回指定name的客户端单例对象
func Instance(config Config) Client {
	redisConfig := SetConfig(config)
	c := Client{
		config: &redisConfig,
	}

	pool := &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				redisConfig.Hostname+":"+redisConfig.Port,
				redis.DialDatabase(redisConfig.Database),
				redis.DialPassword(redisConfig.Password),
				redis.DialUseTLS(redisConfig.UseTLS),
			)
		},
	}

	//用完了就关闭这个链接
	defer pool.Close()

	//从pool连接池里取出一个链接
	conn := pool.Get()

	c.conn = conn

	auth(c) // 身份校验

	c.Db = &Rdb{
		conn: conn,
	}
	c.Scan = &Rscan{
		conn: conn,
	}
	c.String = &Rstring{
		conn: conn,
	}
	c.Expire = &Rexpire{
		conn: conn,
	}
	c.Hash = &Rhash{
		conn: conn,
	}
	c.Set = &Rset{
		conn: conn,
	}
	c.Zset = &Rzset{
		conn: conn,
	}
	c.List = &Rlist{
		conn: conn,
	}
	c.Geo = &Rgeo{
		conn: conn,
	}
	c.Transaction = &Rtransaction{
		conn: conn,
	}
	c.Hyper = &Rhyper{
		conn: conn,
	}
	c.Bit = &Rbit{
		conn: conn,
	}

	return c
}

// 身份校验
func auth(c Client) {
	//链接redis服务，进行身份验证
	if c.config.Username != "" {
		_, err := c.conn.Do("AUTH", c.config.Username, c.config.Password)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := c.conn.Do("AUTH", c.config.Password)
		if err != nil {
			panic(err)
		}
	}
}

// 通用Do方法保留
func (c Client) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.conn.Do(commandName, args...)
}
