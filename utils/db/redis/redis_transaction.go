// Redis transaction 事务管理
package redis

import "github.com/gomodule/redigo/redis"

type Rtransaction struct {
	conn redis.Conn
}

// UNWATCH 解除所有key的监视
// return: 始终返回"OK"
// link: https://redis.io/commands/unwatch/
func (t *Rtransaction) Unwatch() (string, error) {
	return redis.String(t.conn.Do("UNWATCH"))
}

// WATCH 监视所有给定的key
// keys: []string 需要监视的key
// return: 始终返回"OK"
// link: https://redis.io/commands/watch/
func (t *Rtransaction) Watch(keys []string) (string, error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.String(t.conn.Do("WATCH", args...))
}

// DISCARD 刷新事务中所有先前排队的命令并将连接状态恢复为正常。
// return: 始终返回"OK"
// link: https://redis.io/commands/discard/
func (t *Rtransaction) Discard() (string, error) {
	return redis.String(t.conn.Do("DISCARD"))
}

// MULTI 标记事务块的开始。后续命令将使用 EXEC 排队等待原子执行
// return: 始终返回"OK"
// link: https://redis.io/commands/multi/
func (t *Rtransaction) Multi() (string, error) {
	return redis.String(t.conn.Do("MULTI"))
}

// EXEC 执行事务中所有先前排队的命令并将连接状态恢复为正常。 使用 WATCH 时，EXEC 将仅在监视的键未被修改时才执行命令，从而允许检查和设置机制。
// return:
//     1.返回数组：事务块内所有命令的返回值，按命令执行的先后顺序排列。
//     2. 当操作被打断时，返回空值 nil
// 特别说明：因此命令返回值多种多样，故不在此做类型转换，根据业务场景自行转换处理
// link: https://redis.io/commands/exec/
func (t *Rtransaction) Exec() (interface{}, error) {
	return t.conn.Do("EXEC")
}
