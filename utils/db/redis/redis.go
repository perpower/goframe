// redis功能封装--单例模式
// 特别说明：该功能包要求redis版本使用7.0.0以上，引入了redis7.0的一些新特性，否则有些方法会不支持
package redis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	config *Config
	conn   redis.Conn
	Scan   *Rscan
	String *Rstring
	Expire *Rexpire
	Hash   *Rhash
	Set    *Rset
	Zset   *Rzset
	List   *Rlist
	Geo    *Rgeo
}

var (
	once               sync.Once
	defaultFlushdbMode = "SYNC" // FLUSHDB 默认清空模式
	defaultCursor      = 0      // SCAN默认起始游标
	defaultScanNum     = 10     // 单次scan迭代数量
)

// 返回指定name的客户端单例对象
func Instance(config Config) *Client {
	redisConfig := SetConfig(config)
	c := &Client{
		config: &redisConfig,
	}
	once.Do(func() {
		pool := &redis.Pool{
			MaxIdle:     redisConfig.MaxIdle,
			MaxActive:   redisConfig.MaxActive,
			IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					redisConfig.Hostname+":"+redisConfig.Port,
					redis.DialDatabase(redisConfig.Database),
					redis.DialPassword(redisConfig.Password),
					redis.DialUseTLS(redisConfig.UseTLS),
				)
			},
		}

		//用完了就关闭这个链接
		defer pool.Close()

		//从pool连接池里取出一个链接
		conn := pool.Get()

		c.conn = conn
		auth(c)
		initGroup(c, conn)
	})
	return c
}

// 初始化
func initGroup(c *Client, conn redis.Conn) *Client {
	c.Scan = &Rscan{
		conn: conn,
	}
	c.String = &Rstring{
		conn: conn,
	}
	c.Expire = &Rexpire{
		conn: conn,
	}
	c.Hash = &Rhash{
		conn: conn,
	}
	c.Set = &Rset{
		conn: conn,
	}
	c.Zset = &Rzset{
		conn: conn,
	}
	c.List = &Rlist{
		conn: conn,
	}
	c.Geo = &Rgeo{
		conn: conn,
	}

	return c
}

// 身份校验
func auth(c *Client) {
	//链接redis服务，进行身份验证
	if c.config.Username != "" {
		_, err := c.conn.Do("AUTH", c.config.Username, c.config.Password)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := c.conn.Do("AUTH", c.config.Password)
		if err != nil {
			panic(err)
		}
	}
}

// 通用Do方法保留
func (c *Client) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.conn.Do(commandName, args...)
}

// 返回指定key的秒级剩余生存时间
// key: string 键名
// return:
//
//	   reply: int
//				 TTL in seconds
//			     -1: key存在但未设置过过期时间
//			     -2: key不存在
func (c *Client) Ttl(key string) (reply int, err error) {
	return redis.Int(c.Do("TTL", key))
}

// 返回指定的键存在的数量
// keys: []string 键名数组
// return: reply int
func (c *Client) Exists(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.Do("EXISTS", args...))
}

// 删除当前数据库中的所有键，次命令永不失败，慎用！
// mode: string 刷新模式  SYNC(同步)|ASYNC(异步)
func (c *Client) FlushDb(mode string) {
	if mode == "" {
		mode = defaultFlushdbMode
	}
	c.Do("FLUSHDB", mode)
}

// 删除所有数据库中的所有键，此命令永不失败，慎用！
// mode: string 刷新模式  SYNC(同步)|ASYNC(异步)
func (c *Client) FlushAll(mode string) {
	if mode == "" {
		mode = defaultFlushdbMode
	}
	c.Do("FLUSHDB", mode)
}

// 移除指定的key，同步删除，阻塞式，删除小体量简单数据时推荐优先使用该方式
// keys: []string 键名数组
// return: reply int 成功移除的键的数量
func (c *Client) Del(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.Do("DEL", args...))
}

// 移除指定的key，异步删除，非阻塞式，删除大数据时推荐优先使用该方式，以免阻塞其他短io作业
// keys: []string 键名数组
// return: reply int 成功移除的键的数量
func (c *Client) Unlink(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.Do("UNLINK", args...))
}

// 将存储在源键中的值复制到目标键, redis版本不低于6.2.0可用
// source: string 源键
// destination: string 目标键
// return:
//
//	  reply: int
//				0: 拷贝失败
//				1: 拷贝成功
func (c *Client) Copy(source, destination string) (reply int, err error) {
	return redis.Int(c.Do("COPY", source, destination))
}

// 将指定key从当前数据库中移动到目标数据库中
// key: string 键名
// destinationDb: int 目标数据库
// return:
//
//	  reply: int
//				0: 移动失败
//				1: 移动成功
func (c *Client) Move(key string, destinationDb int) (reply int, err error) {
	return redis.Int(c.Do("MOVE", key, destinationDb))
}

// 将指定key重命名，如果newKey是已经存在的key,则实际上会先执行隐式Del操作，
// 所以如果这个newKey包含比较大的值，则可能引起高延迟。
// key: string 指定键名
// newKey: string 新的键名
// return:
//
//	reply: string
func (c *Client) Rename(key, newKey string) (reply string, err error) {
	return redis.String(c.Do("RENAME", key, newKey))
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
func (c *Client) RenameNx(key, newKey string) (reply string, err error) {
	return redis.String(c.Do("RENAMENX", key, newKey))
}

// 返回指定key的类型
// key: 指定键名
// return: typeName string  可能的类型：string、list、set、zset、hash、stream
func (c *Client) Type(key string) (typeName string, err error) {
	return redis.String(c.Do("TYPE", key))
}

// 更改指定键的最后访问时间。如果键不存在，则忽略
// keys: []string 键名数组
// return: reply int 成功修改的键的数量
func (c *Client) Touch(keys []string) (reply int, err error) {
	args := make([]interface{}, 0)
	for _, v := range keys {
		args = append(args, v)
	}
	return redis.Int(c.Do("TOUCH", args...))
}

// DBSIZE, 获取当前选中数据库中的所有键的数量
// return: reply int
func (c *Client) Dbsize() (int, error) {
	return redis.Int(c.Do("DBSIZE"))
}
