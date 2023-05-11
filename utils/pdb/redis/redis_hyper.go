// Redis HyperLogLog
package redis

import "github.com/gomodule/redigo/redis"

type Rhyper struct {
	conn redis.Conn
}

// PFADD 将任意数量的元素添加到指定的 HyperLogLog 里面
// 如果命令执行时给定的键不存在， 那么程序将先创建一个空的 HyperLogLog 结构， 然后再执行命令。
// key: string 给定的key
// elements: []string 元素集
// return:
//
//	  reply: int
//				0: 未发生任何修改
//				1: 如果 HyperLogLog 的内部储存被修改了
// link: https://redis.io/commands/pfadd/
func (h *Rhyper) Pfadd(key string, elements []string) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range elements {
		args = append(args, v)
	}
	return redis.Int(h.conn.Do("PFADD", args...))
}

// PFCOUNT 返回给定key的基数，注意: 基数并不是精确值， 而是一个带有 0.81% 标准错误（standard error）的近似值
// 1. 命令作用于单个键时， 返回储存在给定键的 HyperLogLog 的近似基数， 如果键不存在， 那么返回 0
// 2. 命令作用于多个键时， 返回所有给定 HyperLogLog 的并集的近似基数， 这个近似基数是通过将所有给定 HyperLogLog
// 合并至一个临时 HyperLogLog 来计算得出的。
// keys: []string
// return: int
// link: https://redis.io/commands/pfcount/
func (h *Rhyper) Pfcount(keys []string) (int, error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(h.conn.Do("PFCOUNT", args...))
}

// PFMERGE 将多个 HyperLogLog 合并为一个 HyperLogLog ， 合并后的 HyperLogLog 的基数接近于所有输入 HyperLogLog 的可见集合的并集。
// 如果destkey不存在，则会先创建一个空的 HyperLogLog
// destkey: string 目标键名
// sourcekeys: []string 源键
// return: string  此命令只返回“OK”
// link: https://redis.io/commands/pfmerge/
func (h *Rhyper) Pfmerge(destkey string, sourcekeys []string) (string, error) {
	args := make([]interface{}, 0)
	args = append(args, destkey)
	for _, v := range sourcekeys {
		args = append(args, v)
	}
	return redis.String(h.conn.Do("PFMERGE", args...))
}
