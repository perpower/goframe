// Redis 位图
package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/perpower/goframe/funcs/normal"
)

type Rbit struct {
	conn redis.Conn
}

// SETBIT 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit), 位的设置或清除取决于 value 参数，
// 可以是 0 也可以是 1, 当 key 不存在时，自动生成一个新的字符串值。
// key: string 指定键名
// offset: int 偏移量
// value: int (0 | 1)
// return: int 返回指定偏移量原来储存的位
// link: https://redis.io/commands/setbit/
func (b *Rbit) Setbit(key string, offset, value int) (int, error) {
	return redis.Int(b.conn.Do("SETBIT", key, offset, value))
}

// GETBIT 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。当 offset 比字符串值的长度大，或者 key 不存在时，返回 0
// key: string 指定键名
// offset: int 偏移量
// return: int
// link: https://redis.io/commands/getbit/
func (b *Rbit) Getbit(key string, offset int) (int, error) {
	return redis.Int(b.conn.Do("GETBIT", key, offset))
}

// BITCOUNT 计算给定key存储字符串中，被设置为 1 的比特位的数量，不存在的key会被当做空字符串处理，返回0
// key: string
// start: int 起始位置，可以为负数
// end: int 结束位置，可以为负数
// return: int
// link: https://redis.io/commands/bitcount/
func (b *Rbit) Bitcount(key string, start, end int) (int, error) {
	return redis.Int(b.conn.Do("BITCOUNT", key, start, end))
}

// BITCOUNT 计算给定key存储字符串中，被设置为 1 的比特位的数量，不存在的key会被当做空字符串处理，返回0
// since: 7.0.0
// key: string
// start: int 起始位置，可以为负数
// end: int 结束位置，可以为负数
// indexType: string 索引方式 取值：BYTE | BIT
// return: int
// link: https://redis.io/commands/bitcount/
func (b *Rbit) BitcountIndex(key string, start, end int, indexType string) (int, error) {
	if !normal.InArray(indexType, []string{"BIT", "BYTE"}) {
		indexType = "BYTE"
	}
	return redis.Int(b.conn.Do("BITCOUNT", key, start, end, indexType))
}

// BITPOS 返回位图中第一个值为 指定bit 的二进制位的位置。
// key: string
// bit: int 指定bit位 (0 | 1)
// start: int 起始位置 可以为负数
// end: int 结束位置 可以为负数
// return: int
// link: https://redis.io/commands/bitpos/
func (b *Rbit) Bitpos(key string, bit, start, end int) (int, error) {
	return redis.Int(b.conn.Do("BITPOS", key, start, end))
}

// BITPOS 返回位图中第一个值为 指定bit 的二进制位的位置。
// key: string
// bit: int 指定bit位 (0 | 1)
// start: int 起始位置 可以为负数
// end: int 结束位置 可以为负数
// indexType: string 索引方式 取值：BYTE | BIT
// return: int
// link: https://redis.io/commands/bitpos/
func (b *Rbit) BitposIndex(key string, bit, start, end int, indexType string) (int, error) {
	if !normal.InArray(indexType, []string{"BIT", "BYTE"}) {
		indexType = "BYTE"
	}
	return redis.Int(b.conn.Do("BITPOS", key, start, end, indexType))
}

// BITOP 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
// operation: string 操作类型，取值： AND | OR | XOR | NOT
// destkey: string 结果保存目标键
//
//	AND: 对一个或多个 key 求逻辑并，并将结果保存到 destkey
//	OR: 对一个或多个 key 求逻辑或，并将结果保存到 destkey
//	XOR: 对一个或多个 key 求逻辑异或，并将结果保存到 destkey
//	NOT: 对给定 key 求逻辑非，并将结果保存到 destkey
//
// keys: []string 指定键集
// return: int 保存到 destkey 的字符串的长度，和输入 key 中最长的字符串长度相等。
// link: https://redis.io/commands/bitop/
func (b *Rbit) Bitop(operation, destkey string, keys []string) (int, error) {
	if !normal.InArray(operation, []string{"AND", "OR", "XOR", "NOT"}) {
		operation = "BYTE"
	}

	args := make([]interface{}, 0)
	args = append(args, operation, destkey)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(b.conn.Do("BITOP", args...))
}
