// 接口限流中间件，根据实际需求可注册为全局|指定路由组|指定路由中间件
package middleware

import (
	"time"

	"github.com/perpower/goframe/funcs"
	"github.com/perpower/goframe/pconstants"
	"github.com/perpower/goframe/structs"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

const (
	FillInterval       = 1 * time.Second //默认填充速率
	Capacity     int64 = 200             //默认令牌桶容量
	Quantum      int64 = 10              //默认单次填充令牌数量
)

// fillInterval: 令牌填充间隔
// cap: 令牌桶容量
func RateLimiter(fillInterval time.Duration, cap int64) gin.HandlerFunc {
	if cap < 1 {
		cap = Capacity
	}
	//创建一个令牌桶
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, Quantum)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回系统繁忙提示
		if bucket.TakeAvailable(1) < 1 {
			funcs.OutJson(c, structs.Outjson{
				Code: pconstants.ERROR_3054,
				Msg:  "系统繁忙,请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
