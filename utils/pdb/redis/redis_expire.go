// 生命周期相关
package redis

import "github.com/gomodule/redigo/redis"

type Rexpire struct {
	conn redis.Conn
}

// 设置指定key的生命周期
// key: string 键名
// seconds: int 生存时间 单位s
// options: string NX|XX|GT|LT, 对应值的含义查阅https://redis.io/commands/expire/
// return:
//
//	  reply: int
//				0: 设置失败
//				1: 设置成功
func (c *Rexpire) Expire(key string, seconds int, options string) (reply int, err error) {
	if options == "" {
		return redis.Int(c.conn.Do("EXPIRE", key, seconds))
	}
	return redis.Int(c.conn.Do("EXPIRE", key, seconds, options))
}

// 设置指定key的过期时间点
// key: string 键名
// timestamp: int 秒级时间戳
// options: string NX|XX|GT|LT, 对应值的含义查阅https://redis.io/commands/expireat/
// return:
//
//	  reply: int
//				0: 设置失败
//				1: 设置成功
func (c *Rexpire) ExpireAt(key string, timestamp int, options string) (reply int, err error) {
	if options == "" {
		return redis.Int(c.conn.Do("EXPIREAT", key, timestamp))
	}
	return redis.Int(c.conn.Do("EXPIREAT", key, timestamp, options))
}

// 返回指定key的失效时间秒级时间戳，7.0.0版本以上可用
// key: string 键名
// return:
//
//	   reply: int
//				 Unix timestamp in seconds
//			     -1: key存在但未设置过过期时间
//			     -2: key不存在
func (c *Rexpire) ExpireTime(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("EXPIRETIME", key))
}

// 设置指定key的生命周期
// key: string 键名
// milliseconds: int 生存时间 单位ms
// options: string NX|XX|GT|LT, 7.0.0版本以上支持， 对应值的含义查阅https://redis.io/commands/pexpire/
// return:
//
//	  reply: int
//				0: 设置失败
//				1: 设置成功
func (c *Rexpire) Pexpire(key string, milliseconds int, options string) (reply int, err error) {
	if options == "" {
		return redis.Int(c.conn.Do("EXPIRE", key, milliseconds))
	}
	return redis.Int(c.conn.Do("EXPIRE", key, milliseconds, options))
}

// 设置指定key的过期时间点
// key: string 键名
// miltimestamp: int64 毫秒级时间戳
// options: string NX|XX|GT|LT, 7.0.0版本以上支持，对应值的含义查阅https://redis.io/commands/pexpireat/
// return:
//
//	  reply: int
//				0: 设置失败
//				1: 设置成功
func (c *Rexpire) PexpireAt(key string, miltimestamp int64, options string) (reply int, err error) {
	if options == "" {
		return redis.Int(c.conn.Do("EXPIREAT", key, miltimestamp))
	}
	return redis.Int(c.conn.Do("EXPIREAT", key, miltimestamp, options))
}

// 返回指定key的失效时间秒级时间戳，7.0.0版本以上支持
// key: string 键名
// return:
//
//	   reply: int64
//				 Unix timestamp in seconds
//			     -1: key存在但未设置过过期时间
//			     -2: key不存在
func (c *Rexpire) PexpireTime(key string) (reply int64, err error) {
	return redis.Int64(c.conn.Do("EXPIRETIME", key))
}

// 移除指定key的过期时间，使其变为永不过期
// key: string 键名
// return:
// 		reply: int
//				0: 移除失败
//				1: 移除成功
func (c *Rexpire) Persist(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("PERSIST", key))
}

// 返回指定key的秒级剩余生存时间
// key: string 键名
// return:
//
//	   reply: int
//				 TTL in seconds
//			     -1: key存在但未设置过过期时间
//			     -2: key不存在
func (c *Rexpire) Ttl(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("TTL", key))
}
