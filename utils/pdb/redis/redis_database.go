// Redis Database相关操作
package redis

import "github.com/gomodule/redigo/redis"

type Rdb struct {
	conn redis.Conn
}

// 返回指定的键存在的数量
// keys: []string 键名数组
// return: reply int
func (c *Rdb) Exists(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("EXISTS", args...))
}

// 删除当前数据库中的所有键，次命令永不失败，慎用！
// mode: string 刷新模式  SYNC(同步)|ASYNC(异步)
func (c *Rdb) FlushDb(mode string) {
	if mode == "" {
		mode = defaultFlushdbMode
	}
	c.conn.Do("FLUSHDB", mode)
}

// 删除所有数据库中的所有键，此命令永不失败，慎用！
// mode: string 刷新模式  SYNC(同步)|ASYNC(异步)
func (c *Rdb) FlushAll(mode string) {
	if mode == "" {
		mode = defaultFlushdbMode
	}
	c.conn.Do("FLUSHDB", mode)
}

// 移除指定的key，同步删除，阻塞式，删除小体量简单数据时推荐优先使用该方式
// keys: []string 键名数组
// return: reply int 成功移除的键的数量
func (c *Rdb) Del(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("DEL", args...))
}

// 移除指定的key，异步删除，非阻塞式，删除大数据时推荐优先使用该方式，以免阻塞其他短io作业
// keys: []string 键名数组
// return: reply int 成功移除的键的数量
func (c *Rdb) Unlink(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("UNLINK", args...))
}

// 将存储在源键中的值复制到目标键, redis版本不低于6.2.0可用
// source: string 源键
// destination: string 目标键
// return:
//
//	  reply: int
//				0: 拷贝失败
//				1: 拷贝成功
func (c *Rdb) Copy(source, destination string) (reply int, err error) {
	return redis.Int(c.conn.Do("COPY", source, destination))
}

// 将指定key从当前数据库中移动到目标数据库中
// key: string 键名
// destinationDb: int 目标数据库
// return:
//
//	  reply: int
//				0: 移动失败
//				1: 移动成功
func (c *Rdb) Move(key string, destinationDb int) (reply int, err error) {
	return redis.Int(c.conn.Do("MOVE", key, destinationDb))
}

// 将指定key重命名，如果newKey是已经存在的key,则实际上会先执行隐式Del操作，
// 所以如果这个newKey包含比较大的值，则可能引起高延迟。
// key: string 指定键名
// newKey: string 新的键名
// return:
//
//	reply: string
func (c *Rdb) Rename(key, newKey string) (reply string, err error) {
	return redis.String(c.conn.Do("RENAME", key, newKey))
}

// 将指定key重命名，仅当newKey不存在时才执行，如果指定key不存在，则会返回错误
// 所以如果这个newKey包含比较大的值，则可能引起高延迟。
// key: string 指定键名
// newKey: string 新的键名
// return:
//
//	  reply: int
//				0: newKey当前已经存在
//				1: 重命名成功
func (c *Rdb) RenameNx(key, newKey string) (reply string, err error) {
	return redis.String(c.conn.Do("RENAMENX", key, newKey))
}

// 返回指定key的类型
// key: 指定键名
// return: reply string  可能的类型：string、list、set、zset、hash、stream
func (c *Rdb) Type(key string) (reply string, err error) {
	return redis.String(c.conn.Do("TYPE", key))
}

// 更改指定键的最后访问时间。如果键不存在，则忽略
// keys: []string 键名数组
// return: reply int 成功修改的键的数量
func (c *Rdb) Touch(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("TOUCH", args...))
}

// DBSIZE 获取当前选中数据库中的所有键的数量
// return: reply int
func (c *Rdb) Dbsize() (int, error) {
	return redis.Int(c.conn.Do("DBSIZE"))
}

// SELECT 切换数据库
// return: reply string
func (c *Rdb) Select(index int) (string, error) {
	return redis.String(c.conn.Do("SELECT", index))
}
