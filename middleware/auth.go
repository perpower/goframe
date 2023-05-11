package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/perpower/goframe/funcs/normal"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/utils/perrors"

	"github.com/gin-gonic/gin"
)

// SignHandle 接口验签
// signExpire: time.Duration 验签有效期
// signKey: string 签名秘钥
func SignHandle(signExpire time.Duration, signKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nowtime = ptime.TimestampMilli()
		sign := c.GetHeader("Sign")
		timestamp := c.GetHeader("Timestamp")
		timestamp_int64, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || timestamp_int64 <= 0 {
			timestamp_int64 = 0
		}

		//定义验签参数map集合
		signParams := map[string]string{
			"nonce":     c.GetHeader("Nonce"),
			"timestamp": c.GetHeader("Timestamp"),
		}

		for key, value := range signParams {
			if value == "" {
				c.Abort()
				c.Error(perrors.Newf(perrors.ERROR_1002.Code, "签名错误，参数 %s 不能为空", nil, key))
				return
			}
		}

		//判断签名时效
		if time.Duration(nowtime-timestamp_int64) > signExpire {
			c.Abort()
			c.Error(&perrors.ERROR_1001)
			return
		}

		//组装签名字符串
		signStr := "nonce=" + signParams["nonce"] + "&timestamp=" + signParams["timestamp"] + "&secretkey=" + signKey

		//md5加密
		signStrmd5 := md5.Sum([]byte(signStr))

		//字符转大写
		signStr = fmt.Sprintf("%X", signStrmd5)

		encryptStr := base64.StdEncoding.EncodeToString(normal.String2Bytes(signStr))

		if encryptStr != sign {
			c.Abort()
			c.Error(&perrors.ERROR_1002)
			return
		}

		c.Next()
	}
}
