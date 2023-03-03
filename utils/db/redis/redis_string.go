// 字符串
package redis

import (
	"github.com/perpower/goframe/funcs/convert"
	"github.com/perpower/goframe/funcs/normal"

	"github.com/gomodule/redigo/redis"
)

type Rstring struct {
	conn redis.Conn
}

// SET--全参数功能方法
// key: string 键名
// value: string 值
// condition: string 指定条件，取值范围: NX | XX
// expireType: string 设置生命周期格式， 取值范围：EX | PX | EXAT | PXAT | KEEPTTL
// 参数对照说明可查阅：https://redis.io/commands/set/
func (c *Rstring) Set(key, value, condition, expireType string, timeout int64) (reply string, err error) {
	if condition == "" || !normal.InArray(condition, []string{"NX", "XX"}) {
		switch expireType {
		case "EX|PX|EXAT|PXAT":
			return redis.String(c.conn.Do("SET", key, value, expireType, timeout))
		case "KEEPTTL":
			return redis.String(c.conn.Do("SET", key, value, "KEEPTTL"))
		default:
			return redis.String(c.conn.Do("SET", key, value))
		}
	} else {
		switch expireType {
		case "EX|PX|EXAT|PXAT":
			return redis.String(c.conn.Do("SET", key, value, condition, expireType, timeout))
		case "KEEPTTL":
			return redis.String(c.conn.Do("SET", key, value, condition, "KEEPTTL"))
		default:
			return redis.String(c.conn.Do("SET", key, value, condition))
		}
	}
}

// SET 同时返回旧的值--全参数功能方法
// key: string 键名
// value: string 值
// condition: string 指定条件，取值范围: NX | XX
// expireType: string 设置生命周期格式， 取值范围：EX | PX | EXAT | PXAT | KEEPTTL
// 参数对照说明可查阅：https://redis.io/commands/set/
func (c *Rstring) SetGet(key, value, condition, expireType string, timeout int64) (reply string, err error) {
	if condition == "" || !normal.InArray(condition, []string{"NX", "XX"}) {
		switch expireType {
		case "EX|PX|EXAT|PXAT":
			return redis.String(c.conn.Do("SET", key, value, "GET", expireType, timeout))
		case "KEEPTTL":
			return redis.String(c.conn.Do("SET", key, value, "GET", "KEEPTTL"))
		default:
			return redis.String(c.conn.Do("SET", key, value, "GET"))
		}
	} else {
		switch expireType {
		case "EX|PX|EXAT|PXAT":
			return redis.String(c.conn.Do("SET", key, value, condition, "GET", expireType, timeout))
		case "KEEPTTL":
			return redis.String(c.conn.Do("SET", key, value, condition, "GET", "KEEPTTL"))
		default:
			return redis.String(c.conn.Do("SET", key, value, condition, "GET"))
		}
	}
}

// SETEX，设置key并同时设置秒级生命周期
// key: string 键名
// seconds: int 秒级时间
// value: string 值
// link: https://redis.io/commands/setex/
func (c *Rstring) SetEx(key, value string, seconds int) (reply string, err error) {
	return redis.String(c.conn.Do("SETEX", key, seconds, value))
}

// SETNX, 仅当key无值时执行操作
// key: string 键名
// seconds: int 秒级时间
// value: string 值
// return:
//
//	reply: int
//			0: 设置失败
//			1: 设置成功
//
// link: https://redis.io/commands/setnx/
func (c *Rstring) SetNx(key, value string) (reply int, err error) {
	return redis.Int(c.conn.Do("SETNX", key, value))
}

// PSETEX，设置key并同时设置毫秒级生命周期
// key: string 键名
// milliseconds: int 毫秒级时间
// value: string 值
// link: https://redis.io/commands/psetex/
func (c *Rstring) PsetEx(key, value string, milliseconds int) (reply string, err error) {
	return redis.String(c.conn.Do("PSETEX", key, milliseconds, value))
}

// SETRANGE，从指定偏移量位置开始，覆盖value的整个长度
// key: string 键名
// offset: int 偏移量 可以设置的最大偏移量为 2^29 -1 (536870911)
// value: string 值
// return: reply int 修改之后的字符串长度
// link: https://redis.io/commands/setrange/
func (c *Rstring) SetRange(key, value string, offset int) (reply int, err error) {
	return redis.Int(c.conn.Do("SETRANGE", key, offset, value))
}

// STRLEN，获取指定key存储的字符串值的长度
// key: string 键名
// return: reply int 长度
// link: https://redis.io/commands/strlen/
func (c *Rstring) Strlen(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("STRLEN", key))
}

// MSET,同时为多个键设置值,此操作从不失败，如果某个给定键已经存在，那么 MSET 将使用新值去覆盖旧值
// keyValues: [][2]string
// link: https://redis.io/commands/mset/
func (c *Rstring) Mset(keyValues [][2]string) (reply string, err error) {
	res := make([]interface{}, 0)
	for _, v := range keyValues {
		res = append(res, v[0], v[1])
	}
	return redis.String(c.conn.Do("MSET", res...))
}

// MSETNX,同时为多个键设置值,只要有一个key存在，此命令都不会进行任何操作
// keyValues: [][2]string
// return:
//
//	reply: int
//			0: 设置失败
//			1: 设置成功
//
// link: https://redis.io/commands/msetnx/
func (c *Rstring) MsetNx(keyValues [][2]string) (reply int, err error) {
	res := make([]interface{}, 0)
	for _, v := range keyValues {
		res = append(res, v[0], v[1])
	}
	return redis.Int(c.conn.Do("MSETNX", res...))
}

// GET
// key: string 键名
// return: reply string
// link: https://redis.io/commands/get/
func (c *Rstring) Get(key string) (reply string, err error) {
	return redis.String(c.conn.Do("GET", key))
}

// GETDEL, 获取指定键值并在成功时删除该键
// key: string 键名
// return: reply string
// link: https://redis.io/commands/getdel/
func (c *Rstring) GetDel(key string) (reply string, err error) {
	return redis.String(c.conn.Do("GETDEL", key))
}

// GETDEL, 获取指定键值并设置它的生命周期
// key: string 键名
// return: reply string
// expireType: string 设置生命周期格式， 取值范围：EX | PX | EXAT | PXAT | PERSIST
// link：https://redis.io/commands/getex/
func (c *Rstring) GetEx(key, expireType string, timeout int64) (reply string, err error) {
	switch expireType {
	case "EX|PX|EXAT|PXAT":
		return redis.String(c.conn.Do("GETEX", key, expireType, timeout))
	case "PERSIST":
		return redis.String(c.conn.Do("GETEX", key, "PERSIST"))
	default:
		return redis.String(c.conn.Do("GETEX", key))
	}
}

// GETRANGE，截取指定偏移量区间的字符串并返回
// key: string 键名
// start: int 支持负数
// end: int 支持负数
// return: reply string
// link：https://redis.io/commands/getrange/
func (c *Rstring) GetRange(key string, start, end int) (reply string, err error) {
	return redis.String(c.conn.Do("GETRANGE", key, start, end))
}

// GETSET, 将键 key 的值设为 value ， 并返回键 key 在被设置之前的旧值。6.2.0版本以后已弃用该方法
// key: string 键名
// value: string 值
// return: reply string
// link：https://redis.io/commands/getset/
func (c *Rstring) GetSet(key, value string) (reply string, err error) {
	return redis.String(c.conn.Do("GETSET", key, value))
}

// MGET,同时为多个键设置值,此操作从不失败，如果某个给定键不存在，那么对应的值返回""
// keys: ...string 键名
// link: https://redis.io/commands/mget/
func (c *Rstring) Mget(keys []string) (reply []string, err error) {
	res := make([]interface{}, 0)
	for _, v := range keys {
		res = append(res, v)
	}
	arr, err := redis.Values(c.conn.Do("MGET", res...))
	for _, v := range arr {
		reply = append(reply, convert.String(v))
	}

	return reply, err
}

// INCR , 将存储在 key 中的数字递增 1, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// return: reply int 返回操作之后的新值
// link：https://redis.io/commands/incr/
func (c *Rstring) Incr(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("INCR", key))
}

// INCRBY , 将存储在 key 中的数字递增指定增量, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// increment: int 增量值,可以为负值
// return: reply int 返回操作之后的新值
// link：https://redis.io/commands/incrby/
func (c *Rstring) IncrBy(key string, increment int) (reply int, err error) {
	return redis.Int(c.conn.Do("INCRBY", key, increment))
}

// INCRBYFLOAT , 将存储在 key 中的数字递增指定增量, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// increment: float64 增量值,可以设置负数
// return: reply float64 返回操作之后的新值
// link：https://redis.io/commands/incrbyfloat/
func (c *Rstring) IncrByFloat(key string, increment float64) (reply float64, err error) {
	return redis.Float64(c.conn.Do("INCRBYFLOAT", key, increment))
}

// DECR , 将存储在 key 中的数字递减 1, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// return: reply int 返回操作之后的新值
// link：https://redis.io/commands/decr/
func (c *Rstring) Decr(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("DECR", key))
}

// DECRBY , 将存储在 key 中的数字递减指定量, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// decrement: int 递减值,可以为负值
// return: reply int 返回操作之后的新值
// link：https://redis.io/commands/decrby/
func (c *Rstring) DecrBy(key string, decrement int) (reply int, err error) {
	return redis.Int(c.conn.Do("DECRBY", key, decrement))
}

// APPEND, 在指定key尾部追加字符串
// key: string 键名
// value: string 追加字符串
// return: reply int 追加之后的字符串长度
// link：https://redis.io/commands/append/
func (c *Rstring) Append(key, value string) (reply int, err error) {
	return redis.Int(c.conn.Do("APPEND", key, value))
}
