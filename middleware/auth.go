package middleware

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"go-framework/configs"
	"strconv"
	"strings"
	"time"

	"github.com/perpower/goframe/funcs"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/pconstants"
	"github.com/perpower/goframe/structs"

	"github.com/gin-gonic/gin"
)

// 接口验签
func CheckSign() gin.HandlerFunc {
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
				jsonData := structs.Outjson{
					Code: pconstants.ERROR_1002,
					Msg:  "签名错误，参数" + key + "不能为空",
				}
				funcs.OutJson(c, jsonData)
				c.Abort()
				return
			}
		}

		//判断签名时效
		if time.Duration(nowtime-timestamp_int64) > configs.AuthConfig.Expire {
			jsonData := structs.Outjson{
				Code: pconstants.ERROR_1001,
			}
			funcs.OutJson(c, jsonData)
			c.Abort()
			return
		}

		//组装签名字符串
		signStr := "nonce=" + signParams["nonce"] + "&timestamp=" + signParams["timestamp"] + "&secretkey=" + configs.AuthConfig.Key

		//md5加密
		signStrmd5 := md5.Sum([]byte(signStr))

		//字符转大写
		signStr_pre := fmt.Sprintf("%x", signStrmd5)
		signStr = strings.ToUpper(signStr_pre)

		_decodeStr, _ := base64.StdEncoding.DecodeString(sign)
		block, _ := aes.NewCipher([]byte(configs.AuthConfig.Key))
		decSign := funcs.SubByte(_decodeStr, block.BlockSize(), 0)

		iv := funcs.SubByte(_decodeStr, 0, block.BlockSize())
		encryptStr, _ := funcs.AESEncrypt(signStr, configs.AuthConfig.Key, iv)

		if encryptStr != funcs.Bytes2String(decSign) {
			jsonData := structs.Outjson{
				Code: pconstants.ERROR_1002,
			}
			funcs.OutJson(c, jsonData)
			c.Abort()
			return
		}

		c.Next()
	}
}
