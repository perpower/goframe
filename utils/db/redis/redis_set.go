// 集合
package redis

import (
	"github.com/gomodule/redigo/redis"
)

type Rset struct {
	conn redis.Conn
}

// SADD, 往指定集合添加成员，已经存在的跳过, key存在将创建
// key: string 键名
// members: []string 成员
// return: reply int 成功添加的成员数
// link: https://redis.io/commands/sadd/
func (c *Rset) Sadd(key string, members []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("SADD", args...))
}

// SCARD, 返回指定key中存储的成员数
// key: string 键名
// return: reply int
// link: https://redis.io/commands/scard/
func (c *Rset) Scard(key string) (int, error) {
	return redis.Int(c.conn.Do("SCARD", key))
}

// SDIFF, 返回由第一个集合和所有后续集合之间的差异产生的集合成员
// key: string 第一个集合键名
// diffKeys: []string 其他对比的键名数组
// return: reply int
// link: https://redis.io/commands/sdiff/
func (c *Rset) Sdiff(key string, diffKeys []string) ([]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range diffKeys {
		args = append(args, v)
	}
	return redis.Strings(c.conn.Do("SDIFF", args...))
}

// SDIFF, 比对第一个集合和所有后续集合之间的差异产生的集合成员，
// 并将结果存储到目标destination中，如果destination存在则值会被完全覆盖。
// key: string 第一个集合键名
// diffKeys: []string 其他对比的键名数组
// destination: string 目标键名
// return: reply int destination的成员数量
// link: https://redis.io/commands/sdiffstore/
func (c *Rset) SdiffStore(key string, diffKeys []string, destination string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination, key)
	for _, v := range diffKeys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("SDIFFSTORE", args...))
}

// SINTER, 返回所有给定集合的交集成员，其中只要有一个空集合，结果必然为空，不存在的key也视为空集合
// key: []string 键名数组
// return: reply []string
// link: https://redis.io/commands/sinter/
func (c *Rset) Sinter(keys []string) ([]string, error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Strings(c.conn.Do("SINTER", args...))
}

// SINTERCARD, 返回所有给定集合的交集成员数，其中只要有一个空集合，
// 结果必然为0，不存在的key也视为空集合
// since: 7.0.0
// keys: []string 键名数组
// limit: int 限制返回的数量
// return: reply int
// link: https://redis.io/commands/sintercard/
func (c *Rset) SinterCard(keys []string, limit int) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	if limit > 0 {
		args = append(args, "LIMIT", limit)
	}
	return redis.Int(c.conn.Do("SINTERCARD", args...))
}

// SINTERSTORE, 比对所有给定集合的交集成员数，并将结果存储到目标destination中，如果destination存在则值会被覆盖。
// key: []string 其键名数组
// destination: string 目标键名
// return: reply int destination的成员数量
// link: https://redis.io/commands/sinterstore/
func (c *Rset) SinterStore(keys []string, destination string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("SINTERSTORE", args...))
}

// SISMEMBER, 判断指定member是否在集合中
// key: string 键名
// member: string 成员
// return:
//
//	  reply: int
//				0: 不包含，或key不存在
//				1: 包含
//
// link: https://redis.io/commands/sismember/
func (c *Rset) Sismember(key, member string) (int, error) {
	return redis.Int(c.conn.Do("SISMEMBER", key, member))
}

// SMEMBERS, 返回集合中所有成员，大数据量集合慎用，尤其是生产环境
// key: string 键名
// return: reply []string
// link: https://redis.io/commands/smembers/
func (c *Rset) Smembers(key string) ([]string, error) {
	return redis.Strings(c.conn.Do("SMEMBERS", key))
}

// SMISMEMBER, 判断每个给定的member是否在集合中,
// since: 6.2.0
// key: string 键名
// members: []string 成员数组
// return:
//
//	  reply: []int 返回顺序与members顺序一致
//				0: 不包含，或key不存在
//				1: 包含
//
// link: https://redis.io/commands/smismember/
func (c *Rset) Smismember(key string, members []string) ([]int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	return redis.Ints(c.conn.Do("SMISMEMBER", args...))
}

// SMOVE, 判断每个给定的member是否在集合中,如果destination中当前已经存在该member
// 则会从source中移除该member,destination不做任何改变
// source: string 源集合键名
// destination: string 目标集合键名
// members: 指定成员
// return:
//
//	  reply: int
//				0: 转移失败
//				1: 成功转移
//
// link: https://redis.io/commands/smove/
func (c *Rset) Smove(source, destination, member string) (int, error) {
	return redis.Int(c.conn.Do("SMOVE", source, destination, member))
}

// SPOP，从集合中移除并返回一个或多个随机成员
// key: string 键名
// count: int 指定返回成员个数
// return: reply []string
// link: https://redis.io/commands/spop/
func (c *Rset) Spop(key string, count int) ([]string, error) {
	if count < 1 {
		count = 1
	}
	return redis.Strings(c.conn.Do("SPOP", key, count))
}

// SRANDMEMBER，从集合中返回一个或多个随机成员
// key: string 键名
// count: int 指定返回成员个数
//
//	正数：返回值不存在重复，若超过集合最大成员数，则取最小者
//	负数：返回值可能有重复值，返回值数量始终等于count
//
// return: reply []string
// link: https://redis.io/commands/srandmember/
func (c *Rset) Srandmember(key string, count int) ([]string, error) {
	return redis.Strings(c.conn.Do("SRANDMEMBER", key, count))
}

// SREM, 移除指定成员，不存在的成员会被忽略，不存在的集合始终返回0
// key: string 键名
// members: []string 成员数组
// return: reply int 返回成功移除的成员数
// link: https://redis.io/commands/srem/
func (c *Rset) Srem(key string, members []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("SREM", args...))
}

// SUNION，返回由所有给定集合的并集生成的集合的成员
// keys: []string 键名数组
// return: reply []string
// link: https://redis.io/commands/sunion/
func (c *Rset) Sunion(keys []string) ([]string, error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Strings(c.conn.Do("SUNION", args...))
}

// SUNIONSTORE，返回由所有给定集合的并集生成的集合的成员，
// 并存储到目标键destination中，若destination存在，则值会被覆盖
// keys: []string 键名数组
// destination: string 目标键名
// return: reply []string
// link: https://redis.io/commands/sunionstore/
func (c *Rset) SunionStore(keys []string, destination string) ([]string, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Strings(c.conn.Do("SUNIONSTORE", args...))
}

// 根据条件迭代一次指定key中满足条件的成员
// key: string 键名
// cursor: string 游标
// pattern: string 正则表达式
// count: int 单次迭代成员的数量, 存储的数据体量小时，一般返回都是全部结果，该选项并不起作用，这是正常的
// link：https://redis.io/commands/sscan/
func (c *Rset) ScanOnce(key string, cursor int, pattern string, count int) (int, []string, error) {
	if count < 1 {
		count = defaultScanNum
	}
	args := make([]interface{}, 0)
	args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

	res, err := redis.Values(c.conn.Do("SSCAN", args...))
	cur, _ := redis.Int(res[0], nil)
	arr, _ := redis.Strings(res[1], nil)

	return cur, arr, err
}

// 根据条件迭代指定key中所有满足条件的成员
// key: string 键名
// pattern: string 正则表达式
// count: int 单次迭代成员的数量
// return: arr []string
// link：https://redis.io/commands/sscan/
func (c *Rset) ScanAll(key string, pattern string, count int) (arr []string) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("SSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)
			if len(lists) > 0 {
				arr = append(arr, lists...)
			}
		}
		if cursor == defaultCursor {
			break
		}
	}
	return arr
}

// 根据条件迭代指定key中所有满足条件的成员并删除
// key: string 键名
// pattern: string 正则表达式
// count: int 单次迭代成员的数量
// return: nums int 本次迭代删除的成员数量
// link：https://redis.io/commands/sscan/
func (c *Rset) ScanDel(key string, pattern string, count int) (nums int) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	nums = 0
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("SSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)

			//执行删除
			if len(lists) > 0 {
				num, _ := c.Srem(key, lists)
				nums += num
			}
		}

		if cursor == defaultCursor {
			break
		}
	}
	return nums
}
