// 分布式唯一ID生成工具, 需搭配Redis一起使用
// By default, the ID format follows the original Twitter snowflake format.
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  5 Bit DataCenterID |  5 Bit WorkerID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+

// 1. The ID as a whole is a 63 bit integer stored in an int64
// 2. 41 bits are used to store a timestamp with millisecond precision, using a custom epoch.
// 3. 5 bits are used to store a datacenter id - a range from 0 through 31.
// 4. 5 bits are used to store a workerid id - a range from 0 through 31.
// 5. 12 bits are used to store a sequence number - a range from 0 through 4095.

package snowflake

import (
	"fmt"
	"sync"
	"time"

	"github.com/perpower/goframe/funcs/convert"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/utils/db/redis"
)

const (
	epoch             = int64(1675180800000)                           // 设置起始时间(时间戳/毫秒)：2023-02-01 00:00:00，定义之后不能随便修改，否则可能出现相同ID，有效期69年
	timestampBits     = uint(41)                                       // 时间戳占用位数
	datacenteridBits  = uint(5)                                        // 数据中心id所占位数
	workeridBits      = uint(5)                                        // 机器id所占位数
	sequenceBits      = uint(12)                                       // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值    -1的二进制表示是64位全是1
	datacenteridMax   = int64(-1 ^ (-1 << datacenteridBits))           // 支持的最大数据中心id数量
	workeridMax       = int64(-1 ^ (-1 << workeridBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workeridShift     = sequenceBits                                   // 机器id左移位数
	datacenteridShift = sequenceBits + workeridBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workeridBits + datacenteridBits // 时间戳左移位数
	onceNums          = 50000                                          //单次预生成ID数量
	recreatePercent   = 20                                             //预生成百分比阀值
	redisPre          = "snowflake"                                    //	redis key 分组标识
)

type Snowflake struct {
	sync.Mutex   //加锁，防止并发碰撞
	timestamp    int64
	workerid     int64
	datacenterid int64
	sequence     int64
}

var redisClient *redis.Client

// NewSnowflake
// redisObj： *redis.Client 已经实例化的redis链接对象
// datacenterid: int64
// workerid: int64
func NewSnowflake(redisObj *redis.Client, datacenterid, workerid int64) (*Snowflake, error) {
	redisClient = redisObj
	if datacenterid < 0 || datacenterid > datacenteridMax {
		return nil, fmt.Errorf("datacenterid must be between 0 and %d", datacenteridMax-1)
	}
	if workerid < 0 || workerid > workeridMax {
		return nil, fmt.Errorf("workerid must be between 0 and %d", workeridMax-1)
	}
	return &Snowflake{
		timestamp:    0,
		datacenterid: datacenterid,
		workerid:     workerid,
		sequence:     0,
	}, nil
}

// GetMilliStamp 获取当前毫秒时间戳
// return: int64
func getMilliStamp() int64 {
	return ptime.TimestampMilli()
}

// Generate 生成一个唯一ID
// return: string
func (s *Snowflake) Generate() (string, error) {
	exit, _ := redisClient.Db.Exists([]string{redisPre + ":idsList"})
	if exit == 0 {
		s.Produce(onceNums)
	}

	res, err := redisClient.Zset.ZpopMin(redisPre+":idsList", 1)
	if err != nil || len(res) == 0 {
		return "", err
	}

	return res[0][1], err
}

// GenerateBatch 批量生成指定数量的ID
// nums: int 数量
// return: []string 返回ID数组
func (s *Snowflake) GenerateBatch(nums int) (arr []string, err error) {
	count, err := redisClient.Zset.Zcard(redisPre + ":idsList")
	if err != nil {
		return []string{}, err
	}
	if (count <= (recreatePercent*onceNums)/100) || (count < nums) {
		s.Produce(nums)
	}

	res, err := redisClient.Zset.ZpopMin(redisPre+":idsList", nums)
	if err != nil || len(res) == 0 {
		return []string{}, err
	}

	for _, val := range res {
		arr = append(arr, val[1])
	}
	return arr, err
}

// Produce 批量预生成ID，并将结果存储到redis中
// return:
//
//	count: 成功写入Redis集合的数量
//	err: 错误信息
func (s *Snowflake) Produce(nums int) (int, error) {
	if nums <= onceNums { // 如果此次要获取的数量大于设定的onceNums，则直接生成nums数量的ID
		nums = onceNums
	}

	var i int
	scoreElements := make([][2]string, 0)

	s.Lock()
	for i = 0; i < nums; i++ {
		now := getMilliStamp() // 当前毫秒时间戳
		if s.timestamp == now {
			// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
			s.sequence = (s.sequence + 1) & sequenceMask
			if s.sequence == 0 {
				// 如果当前序列超出12bit长度，则需要等待下一毫秒
				// 下一毫秒将使用sequence:0
				for now <= s.timestamp {
					now = getMilliStamp()
				}
			}
		} else {
			// 不同时间戳（精度：毫秒）下直接使用序列号：0
			s.sequence = 0
		}
		t := now - epoch
		if t > timestampMax {
			continue
		}
		s.timestamp = now
		r := convert.String(int64((t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence)))
		scoreElements = append(scoreElements, [2]string{r, r})
	}

	// 插入redis有序集合中
	count, err := redisClient.Zset.Zadd(redisPre+":idsList", scoreElements, "", "", false)

	s.Unlock()
	return count, err
}

// 获取数据中心ID和机器ID
func GetDeviceID(sid int64) (datacenterid, workerid int64) {
	datacenterid = (sid >> datacenteridShift) & datacenteridMax
	workerid = (sid >> workeridShift) & workeridMax
	return
}

// 获取时间戳
func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

// 获取创建ID时的时间戳
func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + epoch
	return
}

// 获取创建ID时的时间字符串(精度：秒)
func GetGenTime(sid int64) (t string) {
	// 需将GetGenTimestamp获取的时间戳/1000转换成秒
	t = time.Unix(GetGenTimestamp(sid)/1000, 0).Format("2006-01-02 15:04:05")
	return
}

// 获取时间戳已使用的占比：范围（0.0 - 1.0）
func GetTimestampStatus() (state float64) {
	state = float64((getMilliStamp() - epoch)) / float64(timestampMax)
	return
}
