// scan迭代，针对全数据库的键，返回的每个元素都是一个数据库键
package redis

import (
	"github.com/gomodule/redigo/redis"
)

type Rscan struct {
	conn redis.Conn
}

// 根据条件迭代一次当前数据库中满足条件的键集
// cursor: string 游标
// pattern: string 正则表达式
// count: int 单次迭代键的数量
// typ: string 类型 6.0.0版本以后支持该参数
// link：https://redis.io/commands/scan/
func (s *Rscan) ScanOnce(cursor int, pattern string, count int, typ string) (int, []string, error) {
	if count < 1 {
		count = defaultScanNum
	}
	args := make([]interface{}, 0)
	if typ == "" {
		args = append(args, cursor, "MATCH", pattern, "COUNT", count)
	} else {
		args = append(args, cursor, "MATCH", pattern, "COUNT", count, "TYPE", typ)
	}
	res, err := redis.Values(s.conn.Do("SCAN", args...))
	cur, _ := redis.Int(res[0], nil)
	arr, _ := redis.Strings(res[1], nil)

	return cur, arr, err
}

// 根据条件迭代当前数据库中所有满足条件的键集
// pattern: string 正则表达式
// count: int 单次迭代键的数量
// typ: string 类型 6.0.0版本以后支持该参数
// return: arr []string
// link：https://redis.io/commands/scan/
func (s *Rscan) ScanAll(pattern string, count int, typ string) (arr []string) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	for {
		args := make([]interface{}, 0)
		if typ == "" {
			args = append(args, cursor, "MATCH", pattern, "COUNT", count)
		} else {
			args = append(args, cursor, "MATCH", pattern, "COUNT", count, "TYPE", typ)
		}

		res, err := redis.Values(s.conn.Do("SCAN", args...))
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

// 根据条件迭代当前数据库中所有满足条件的键集并删除
// pattern: string 正则表达式
// count: int 单次迭代键的数量
// typ: string 类型 6.0.0版本以后支持该参数
// return: nums int 本次迭代删除的键数量
// link：https://redis.io/commands/scan/
func (s *Rscan) ScanDel(pattern string, count int, typ string) (nums int) {
	if count < 1 {
		count = defaultScanNum
	}
	cursor := defaultCursor
	nums = 0
	for {
		args := make([]interface{}, 0)
		if typ == "" {
			args = append(args, cursor, "MATCH", pattern, "COUNT", count)
		} else {
			args = append(args, cursor, "MATCH", pattern, "COUNT", count, "TYPE", typ)
		}

		res, err := redis.Values(s.conn.Do("SCAN", args...))
		if err == nil {
			curs, _ := redis.Int(res[0], nil)
			cursor = curs
			lists, _ := redis.Values(res[1], nil)

			//执行删除
			if len(lists) > 0 {
				num, _ := redis.Int(s.conn.Do("UNLINK", lists...))
				nums += num
			}
		}

		if cursor == defaultCursor {
			break
		}
	}
	return nums
}
