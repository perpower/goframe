package redis

import (
	"github.com/gomodule/redigo/redis"
)

// 根据条件迭代一次指定key中满足条件的成员
// key: string 键名
// cursor: string 游标
// pattern: string 正则表达式
// count: int 单次迭代成员的数量, 存储的数据体量小时，一般返回都是全部结果，该选项并不起作用，这是正常的
// link：https://redis.io/commands/zscan/
func (c *Rzset) ScanOnce(key string, cursor int, pattern string, count int) (int, [][2]string, error) {
	if count < 1 {
		count = defaultScanNum
	}
	args := make([]interface{}, 0)
	args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

	res, err := redis.Values(c.conn.Do("ZSCAN", args...))
	cur, _ := redis.Int(res[0], nil)
	arr, _ := redis.Strings(res[1], nil)

	// 将键值对数组进行拆分
	resArr := make([][2]string, 0)
	var i int
	for i = 0; i < len(arr); i++ {
		resArr = append(resArr, [2]string{arr[i+1], arr[i]})
		i++
	}

	return cur, resArr, err
}

// 根据条件迭代指定key中所有满足条件的成员
// key: string 键名
// pattern: string 正则表达式
// count: int 单次迭代成员的数量
// return: arr []string
// link：https://redis.io/commands/zscan/
func (c *Rzset) ScanAll(key string, pattern string, count int) (arr [][2]string) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	var i int
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("ZSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)
			if len(lists) > 0 {
				// 将键值对数组进行拆分
				for i = 0; i < len(lists); i++ {
					arr = append(arr, [2]string{lists[i+1], lists[i]})
					i++
				}
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
// link：https://redis.io/commands/zscan/
func (c *Rzset) ScanDel(key string, pattern string, count int) (nums int) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	nums = 0
	var i int
	for {
		args := make([]interface{}, 0)
		args = append(args, key, cursor, "MATCH", pattern, "COUNT", count)

		res, err := redis.Values(c.conn.Do("ZSCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Strings(res[1], nil)

			//执行删除
			if len(lists) > 0 {
				members := []string{}
				for i = 0; i < len(lists); i++ {
					if i%2 == 0 {
						members = append(members, lists[i])
					}
				}
				num, _ := c.Zrem(key, members)
				nums += num
			}
		}

		if cursor == defaultCursor {
			break
		}
	}
	return nums
}
