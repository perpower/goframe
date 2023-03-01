// 列表
package redis

import (
	"github.com/perpower/goframe/funcs"

	"github.com/gomodule/redigo/redis"
)

type Rlist struct {
	conn redis.Conn
}

// LINDEX, 返回指定索引位的元素
// key: string 键名
// index: int 索引，可以负数，当为负数时表示从列表尾部开始
// return: reply string
// link: https://redis.io/commands/lindex/
func (c *Rlist) Lindex(key string, index int) (string, error) {
	return redis.String(c.conn.Do("LINDEX", key, index))
}

// LINSERT, 在列表指定位置插入元素
// key: string 键名
// condition: string 指定插入的方式 BEFORE | AFTER
// pivot: string  参考值
// element: string 待插入的元素
// return:
//
//	  reply: int 成功插入元素之后列表的长度
//				0: key不存在
//				-1: 参考值pivot未找到
//
// link: https://redis.io/commands/linsert/
func (c *Rlist) Linsert(key, condition, pivot, element string) (int, error) {
	if condition == "" || !funcs.InArray(condition, []string{"BEFORE", "AFTER"}) {
		panic("参数condition传值不正确,值必须为 BEFORE | AFTER")
	}
	return redis.Int(c.conn.Do("LINSERT", key, condition, pivot, element))
}

// LLEN, 返回指定列表的长度
// key: string 键名
// return: reply int
// link: https://redis.io/commands/llen/
func (c *Rlist) Llen(key string) (int, error) {
	return redis.Int(c.conn.Do("LLEN", key))
}

// LMOVE, 以原子方式返回并删除存储在源中的列表的第一个/最后一个元素（头/尾取决于 wherefrom 参数），
// 并将元素存入目标列表的第一个/最后一个元素（头/尾取决于 whereto 参数）
// 如果源不存在，则返回值nil，不执行任何操作。 如果 source 和 destination 相同，
// 则该操作相当于从列表中删除第一个/最后一个元素并将其作为列表的第一个/最后一个元素推送，
// 因此可以将其视为列表轮换命令（或空操作 如果wherefrom与whereto相同）。
// since: 6.2.0
// source: string 源列表
// destination: string 目标列表
// wherefrom: string 源列表取出元素的位置 取值：LEFT | RIGHT
// whereto: string 目标列表存入元素的位置 取值：LEFT | RIGHT
// return: reply string 返回本次操作的元素
// link: https://redis.io/commands/lmove/
func (c *Rlist) Lmove(source, destination, wherefrom, whereto string) (string, error) {
	if wherefrom == "" || !funcs.InArray(wherefrom, []string{"LEFT", "RIGHT"}) {
		panic("参数wherefrom传值不正确,值必须为 LEFT | RIGHT")
	}
	if whereto == "" || !funcs.InArray(whereto, []string{"LEFT", "RIGHT"}) {
		panic("参数whereto传值不正确,值必须为 LEFT | RIGHT")
	}
	return redis.String(c.conn.Do("LMOVE", source, destination, wherefrom, whereto))
}

// LMPOP, 从提供的键名列表中的第一个非空列表键中移出一个或多个元素。
// since: 7.0.0
// keys: []string 键名数组
// condition: string  取值：LEFT | RIGHT
// count: int 返回元素的数量，取非空列表的长度和count两者的较小者
// return:
// link: https://redis.io/commands/lmpop/
func (c *Rlist) Lmpop(keys []string, condition string, count int) (keyName string, arr []string, err error) {
	if condition == "" || !funcs.InArray(condition, []string{"LEFT", "RIGHT"}) {
		panic("参数condition传值不正确,值必须为 LEFT | RIGHT")
	}
	if count < 1 {
		count = 1
	}
	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	args = append(args, condition, "COUNT", count)
	res, err := redis.Values(c.conn.Do("LMPOP", args...))

	if len(res) == 0 || err != nil {
		return "", []string{}, err
	}

	keyName, _ = redis.String(res[0], nil)
	arr, _ = redis.Strings(res[1], nil)

	return keyName, arr, err
}

// LPOP, 从列表的开头开始删除并返回一个或多个元素
// key: string 键名
// count: int 需移出的数量
// return: reply []string
// link: https://redis.io/commands/lpop/
func (c *Rlist) Lpop(key string, count int) ([]string, error) {
	if count < 1 {
		count = 1
	}

	return redis.Strings(c.conn.Do("LPOP", key, count))
}

// LPOS，根据条件返回元素的在列表中的索引位置
// since: 6.0.6
// key: string 键名
// element: string 待匹配的元素
// rank: int 正数代表从头开始搜索，选项指定要返回的第一个元素的“等级”，以防存在多个匹配项。
// 等级 1 表示返回第一个匹配项，等级 2 表示返回第二个匹配项，依此类推。
// 若值为负数，代表从尾部开始搜索，但是索引值不会受搜索方向影响。
// count: int 指定返回前N个匹配的元素位置，0代表所有匹配的
// maxlen: int MAXLEN 选项告诉命令只将提供的元素与给定的maxlen次数进行比较
// 当使用 MAXLEN 时，可以将 0 指定为最大比较次数，以此告诉命令我们需要无限次比较。
// 这比给出一个非常大的 MAXLEN 选项要好，因为它更通用。
// return: reply []int
// link: https://redis.io/commands/lpos/
func (c *Rlist) Lpos(key, element string, rank, count, maxlen int) ([]int, error) {
	if element == "" {
		return []int{}, nil
	}
	if rank == 0 {
		rank = 1
	}
	return redis.Ints(c.conn.Do("LPOS", key, element, "RANK", rank, "COUNT", count, "MAXLEN", maxlen))
}

// LPUSH, 往列表头部插入多个元素, 若列表key不存在，则先创建
// key: string 键名
// elements: []string 待插入的元素数组
// return: reply int 返回操作执行之后列表的长度
// link: https://redis.io/commands/lpush/
func (c *Rlist) Lpush(key string, elements []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range elements {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("LPUSH", args...))
}

// LPUSHX, 往列表头部插入多个元素, 若列表key不存在，则不执行任何操作
// key: string 键名
// elements: []string 待插入的元素数组
// return: reply int 返回操作执行之后列表的长度
// link: https://redis.io/commands/lpushx/
func (c *Rlist) Lpushx(key string, elements []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range elements {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("LPUSHX", args...))
}

// LRANGE, 返回列表指定偏移量区间的所有元素
// 超出范围的索引不会产生错误。如果开始大于列表的末尾，则返回一个空列表。
// 如果 stop 大于列表的实际末尾，Redis 会将其视为列表的最后一个元素。
// key: string 键名
// start: int 起始偏移量，可以为负数
// stop: int 截止偏移量，可以为负数
// return: reply []string
// link: https://redis.io/commands/lrange/
func (c *Rlist) Lrange(key string, start, stop int) ([]string, error) {
	return redis.Strings(c.conn.Do("LRANGE", key, start, stop))
}

// LREM, 从列表中删除指定次数的元素
// key: string 键名
// count: int 指定删除的元素次数
// count > 0：删除等于从头到尾移动的元素的元素。
// count < 0：删除等于从尾部移动到头部的元素的元素。
// count = 0：移除所有等于element的元素。
// element: string 元素
// return: reply int 返回被删除元素的数量，若key不存在，始终返回0
// link: https://redis.io/commands/lrange/
func (c *Rlist) Lrem(key string, count int, element string) (int, error) {
	return redis.Int(c.conn.Do("LREM", key, count, element))
}

// LSET, 将列表中索引位的元素设置为指定元素
// key: string 键名
// index: int 索引值 可为负数
// element: string 元素值
// return: reply string  成功返回"OK"
// link: https://redis.io/commands/lset/
func (c *Rlist) Lset(key string, index int, element string) (string, error) {
	return redis.String(c.conn.Do("LSET", key, index, element))
}

// LTRIM, 修剪现有列表，使其仅包含指定索引范围的元素。
// 超出范围的索引不会产生错误：如果 start 大于列表的末尾，或者 start > end，结果将是一个空列表（这会导致 key 被删除）。
// 如果 end 大于列表的末尾，Redis 会将其视为列表的最后一个元素。
// key: string 键名
// start: int 起始索引，可以为负数
// stop: int 截止索引，可以为负数
// return: reply string  成功返回"OK"
// link: https://redis.io/commands/ltrim/
func (c *Rlist) Ltrim(key string, start, stop int) (string, error) {
	return redis.String(c.conn.Do("LTRIM", key, start, stop))
}

// RPOP, 从列表的尾部开始删除并返回一个或多个元素
// key: string 键名
// count: int 需移出的数量
// return: reply []string
// link: https://redis.io/commands/rpop/
func (c *Rlist) Rpop(key string, count int) ([]string, error) {
	if count < 1 {
		count = 1
	}

	return redis.Strings(c.conn.Do("RPOP", key, count))
}

// RPUSH, 往列表尾部插入多个元素, 若列表key不存在，则先创建
// key: string 键名
// elements: []string 待插入的元素数组
// return: reply int 返回操作执行之后列表的长度
// link: https://redis.io/commands/rpush/
func (c *Rlist) Rpush(key string, elements []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range elements {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("RPUSH", args...))
}

// RPUSHX, 往列表尾部插入多个元素, 若列表key不存在，则不执行任何操作
// key: string 键名
// elements: []string 待插入的元素数组
// return: reply int 返回操作执行之后列表的长度
// link: https://redis.io/commands/rpushx/
func (c *Rlist) Rpushx(key string, elements []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range elements {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("RPUSHX", args...))
}
