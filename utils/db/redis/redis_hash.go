// 哈希表
package redis

import (
	"github.com/perpower/goframe/funcs/convert"

	"github.com/gomodule/redigo/redis"
)

type Rhash struct {
	conn redis.Conn
}

// HSET,可同时设置多个field=>value
// key: string 键名
// fieldValues: [][2]string{{field, value}, ...}
// return: reply int  成功添加的字段数
// link: https://redis.io/commands/hset/
func (c *Rhash) Hset(key string, fieldValues [][2]string) (reply int, err error) {
	arr := make([]interface{}, 0)
	arr = append(arr, key)
	for _, v := range fieldValues {
		arr = append(arr, v[0], v[1])
	}
	return redis.Int(c.conn.Do("HSET", arr...))
}

// HSETNX, 一次只能设置一个field
// key: string 键名
// field: string 字段名
// value: string 字段值
// return:
//
//	  reply: int
//				0: 设置失败
//				1: 设置成功,当且仅当field不存在时
//
// link: https://redis.io/commands/hsetnx/
func (c *Rhash) HsetNx(key, field, value string) (reply int, err error) {
	return redis.Int(c.conn.Do("HSETNX", key, field, value))
}

// HSTRLEN，返回指定field存储的value值字符串长度
// key: string 键名
// field: string 字段名
// return: reply: int 若field不存在返回0
// link: https://redis.io/commands/hstrlen/
func (c *Rhash) HstrLen(key, field string) (reply int, err error) {
	return redis.Int(c.conn.Do("HSTRLEN", key, field))
}

// HVALS, 返回指定key中存储的所有字段的value， key若不存在返回空数组
// key: string 键名
// return: reply []string
// link: https://redis.io/commands/hvals/
func (c *Rhash) Hvals(key string) (reply []string, err error) {
	return redis.Strings(c.conn.Do("HVALS", key))
}

// HKEYS, 返回指定key中存储的所有field， key若不存在返回空数组
// key: string 键名
// return: reply []string
// link: https://redis.io/commands/hkeys/
func (c *Rhash) Hkeys(key string) (reply []string, err error) {
	return redis.Strings(c.conn.Do("HKEYS", key))
}

// HLEN, 返回指定key中存储的所有field的数量
// key: string 键名
// return: reply int
// link: https://redis.io/commands/hlen/
func (c *Rhash) Hlen(key string) (reply int, err error) {
	return redis.Int(c.conn.Do("HLEN", key))
}

// HDEL, 从key中移除指定fields
// key: string 键名
// fields: []string 字段集
// return: reply int 返回成功移除的field数量
// link: https://redis.io/commands/hdel/
func (c *Rhash) Hdel(key string, fields []string) (reply int, err error) {
	arr := make([]interface{}, 0)
	arr = append(arr, key)
	for _, v := range fields {
		arr = append(arr, v)
	}
	return redis.Int(c.conn.Do("HDEL", arr...))
}

// HEXISTS, 判断key中是否包含指定fields
// key: string 键名
// field: string 字段名
// return:
//
//	  reply: int
//				0: 不包含，或key不存在
//				1: 包含
//
// link: https://redis.io/commands/hexists/
func (c *Rhash) Hexists(key, field string) (reply int, err error) {
	return redis.Int(c.conn.Do("HEXISTS", key, field))
}

// HGET, 获取key中指定field的值
// key: string 键名
// field: string 字段名
// return: reply string
// link: https://redis.io/commands/hget/
func (c *Rhash) Hget(key, field string) (reply string, err error) {
	return redis.String(c.conn.Do("HGET", key, field))
}

// HGETALL, 返回存储在key中的所有字段和值
// key: string 键名
// return: reply map[string]string
// link: https://redis.io/commands/hgetall/
func (c *Rhash) HgetAll(key string) (reply map[string]string, err error) {
	return redis.StringMap(c.conn.Do("HGETALL", key))
}

// HMGET, 返回key中指定field的值
// key: string 键名
// fields: []string 字段集
// return: reply []string
// link: https://redis.io/commands/hmget/
func (c *Rhash) Hmget(key string, fields []string) (reply []string, err error) {
	arr := make([]interface{}, 0)
	arr = append(arr, key)
	for _, v := range fields {
		arr = append(arr, v)
	}
	return redis.Strings(c.conn.Do("HMGET", arr...))
}

// HINCRBY , 将存储在 key 中的指定field值递增指定增量,
// 如果key不存在，则创建，若field不存在，则先创建并将值设为初始0
// key: string 键名
// field: string 字段名
// increment: int 增量值,可以为负值
// return: reply int 返回操作之后的新值
// link：https://redis.io/commands/incrby/
func (c *Rhash) HincrBy(key, field string, increment int) (reply int, err error) {
	return redis.Int(c.conn.Do("HINCRBY", key, field, increment))
}

// HINCRBYFLOAT , 将存储在 key 中的数字递增指定增量, 如果key不存在，则在执行操作之前将其设置为 0
// key: string 键名
// field: string 字段名
// increment: float64 增量值,可以设置负数
// return: reply float64 返回操作之后的新值
// link：https://redis.io/commands/hincrbyfloat/
func (c *Rhash) HincrByFloat(key, field string, increment float64) (reply float64, err error) {
	return redis.Float64(c.conn.Do("HINCRBYFLOAT", key, field, increment))
}

// 根据条件迭代一次指定key中满足条件的键值对
// key: string 键名
// cursor: string 游标
// pattern: string 正则表达式
// count: int 单次迭代键值对的数量, 存储的数据体量小时，一般返回都是全部结果，该选项并不起作用，这是正常的
// link：https://redis.io/commands/hscan/
func (c *Rhash) ScanOnce(key string, cursor int, pattern string, count int) (int, [][2]string, error) {
	if count < 1 {
		count = defaultScanNum
	}
	args := make([]interface{}, 0)
	args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

	res, err := redis.Values(c.conn.Do("HSCAN", args...))
	cur, _ := redis.Int(res[0], nil)
	arr, _ := redis.Strings(res[1], nil)

	// 将键值对数组进行拆分
	resArr := make([][2]string, 0)
	var i int
	for i = 0; i < len(arr); i++ {
		resArr = append(resArr, [2]string{arr[i], arr[i+1]})
		i++
	}

	return cur, resArr, err
}

// 根据条件迭代指定key中所有满足条件的键值对
// key: string 键名
// pattern: string 正则表达式
// count: int 单次迭代键值对的数量
// return: arr [][2]string
// link：https://redis.io/commands/hscan/
func (c *Rhash) ScanAll(key string, pattern string, count int) [][2]string {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	resArr := make([][2]string, 0)
	var i int
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("HSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)
			if len(lists) > 0 {
				// 将键值对数组进行拆分
				for i = 0; i < len(lists); i++ {
					resArr = append(resArr, [2]string{lists[i], lists[i+1]})
					i++
				}
			}
		}
		if cursor == defaultCursor {
			break
		}
	}
	return resArr
}

// 根据条件迭代指定key中所有满足条件的键值对并删除
// key: string 键名
// pattern: string 正则表达式
// count: int 单次迭代键值对的数量
// return: nums int 本次迭代删除的键值对数量
// link：https://redis.io/commands/hscan/
func (c *Rhash) ScanDel(key string, pattern string, count int) (nums int) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	nums = 0
	var i int
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("HSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)

			//执行删除
			if len(lists) > 0 {
				delFields := make([]string, 0)
				for i = 0; i < len(lists); i++ {
					if i%2 == 0 {
						delFields = append(delFields, convert.String(lists[i]))
					}
				}
				num, _ := c.Hdel(key, delFields)
				nums += num
			}
		}

		if cursor == defaultCursor {
			break
		}
	}
	return nums
}
