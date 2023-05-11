package redis

import (
	"github.com/perpower/goframe/funcs/normal"

	"github.com/gomodule/redigo/redis"
)

// ZRANGE, 返回集合中指定范围区间的成员，此方法默认不指定 WITHSCORES参数
// since: 6.2.0
// key: string 键名
// start: string 起始区间
// stop: string 截止区间
// byType: string 范围查询方式， 取值：BYSCORE | BYLEX
// isRev: bool, 成员的顺序默认score从低到高，若需要反向顺序，需指定此参数为true
// limit: []int{offset,count}, LIMIT 参数可用于从匹配元素中获取子范围, 类似 sql的 [offset,count]
// 负数 <count> 返回 <offset> 中的所有元素。 请记住，如果 <offset> 很大，则需要遍历排序集以获取 <offset> 元素，然后才能返回元素
// return:
//
//	reply: []string
//
// link：https://redis.io/commands/zrange/
func (c *Rzset) Zrange(key, start, stop, byType string, isRev bool, limit []int) ([]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key, start, stop)
	if normal.InArray(byType, []string{"BYSCORE", "BYLEX"}) {
		args = append(args, byType)
	}

	if isRev {
		args = append(args, "REV")
	}

	if len(limit) > 0 {
		args = append(args, "LIMIT", limit[0], limit[1])
	}

	return redis.Strings(c.conn.Do("ZRANGE", args...))
}

// ZRANGE, 返回集合中指定范围区间的成员，此方法默认指定 WITHSCORES参数
// since: 6.2.0
// key: string 键名
// start: string 起始区间
// stop: string 截止区间
// byType: string 范围查询方式， 取值：BYSCORE | BYLEX
// isRev: bool, 成员的顺序默认score从低到高，若需要反向顺序，需指定此参数为true
// limit: []int{offset,count}, LIMIT 参数可用于从匹配元素中获取子范围, 类似 sql的 [offset,count]
// 负数 <count> 返回 <offset> 中的所有元素。 请记住，如果 <offset> 很大，则需要遍历排序集以获取 <offset> 元素，然后才能返回元素
// return:
//
//	reply: [][2]string{{score,element}, ...}
//
// link：https://redis.io/commands/zrange/
func (c *Rzset) ZrangeWithScore(key, start, stop, byType string, isRev bool, limit []int) (arr [][2]string, err error) {
	args := make([]interface{}, 0)
	args = append(args, key, start, stop)
	if normal.InArray(byType, []string{"BYSCORE", "BYLEX"}) {
		args = append(args, byType)
	}

	if isRev {
		args = append(args, "REV")
	}

	if len(limit) > 0 {
		args = append(args, "LIMIT", limit[0], limit[1])
	}

	args = append(args, "WITHSCORES")

	res, err := redis.Strings(c.conn.Do("ZRANGE", args...))

	if err != nil {
		return [][2]string{}, err
	}
	if len(res) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(res); i++ {
			arr = append(arr, [2]string{res[i+1], res[i]})
			i++
		}
	}

	return arr, err
}

// ZRANGESTORE, 返回集合中指定范围区间的成员，此方法默认不指定 WITHSCORES参数
// since: 6.2.0
// dst: string 目标键名
// src: string 源键名
// min: string 起始区间
// max: string 截止区间
// byType: string 范围查询方式， 取值：BYSCORE | BYLEX
// isRev: bool, 成员的顺序默认score从低到高，若需要反向顺序，需指定此参数为true
// limit: []int{offset,count}, LIMIT 参数可用于从匹配元素中获取子范围, 类似 sql的 [offset,count]
// 负数 <count> 返回 <offset> 中的所有元素。 请记住，如果 <offset> 很大，则需要遍历排序集以获取 <offset> 元素，然后才能返回元素
// return:
//
//	reply: int  存储到目标集合的成员数
//
// link：https://redis.io/commands/zrangestore/
func (c *Rzset) ZrangeStore(dst, src, min, max, byType string, isRev bool, limit []int) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, dst, src, min, max)
	if normal.InArray(byType, []string{"BYSCORE", "BYLEX"}) {
		args = append(args, byType)
	}

	if isRev {
		args = append(args, "REV")
	}

	if len(limit) > 0 {
		args = append(args, "LIMIT", limit[0], limit[1])
	}

	return redis.Int(c.conn.Do("ZRANGESTORE", args...))
}

// ZREMRANGEBYLEX, 当有序集合中的所有元素都以相同的score插入时，为了强制按字典顺序排序，
// 此命令删除由 min 和 max 指定的字典序范围之间的所有元素。
// key: string 键名
// min: string 起始值
// max: string 结束值
// return:
//
//	reply: int 返回被删除的元素数量
//
// link：https://redis.io/commands/zremrangebylex/
func (c *Rzset) ZremRangeByLex(key, min, max string) (int, error) {
	return redis.Int(c.conn.Do("ZREMRANGEBYLEX", key, min, max))
}

// ZREMRANGEBYRANK, 删除存储在 key 处且排名在 start 和 stop 之间的排序集中的所有元素。
// start 和 stop 都是基于 0 的索引，其中 0 是score最低的元素。 这些索引可以是负数，
// 表示从score最高的元素开始的偏移量。 例如：-1 是score最高的元素，-2 是score第二高的元素，依此类推。
// key: string 键名
// start: int 起始值
// stop: int 结束值
// return:
//
//	reply: int 返回被删除的元素数量
//
// link：https://redis.io/commands/zremrangebyrank/
func (c *Rzset) ZremRangeByRank(key string, min, max int) (int, error) {
	return redis.Int(c.conn.Do("ZREMRANGEBYRANK", key, min, max))
}

// ZREMRANGEBYSCORE, 删除由 min 和 max 指定的score范围之间的所有元素。
// key: string 键名
// min: string 起始值
// max: string 结束值
// return:
//
//	reply: int 返回被删除的元素数量
//
// link：https://redis.io/commands/zremrangebyscore/
func (c *Rzset) ZremRangeByScore(key, min, max string) (int, error) {
	return redis.Int(c.conn.Do("ZREMRANGEBYSCORE", key, min, max))
}
