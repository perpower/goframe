// 短信发送工具包
package psms

type Sms struct {
	TencentSms *tencentSms
}

// Instance
// conf: interface{} 短信发送方式配置
func Instance(conf interface{}) (s Sms, typ interface{}) {
	switch confType := conf.(type) { // 考虑到switch类型断言的问题，将结果分配给一个变量，否则可能会触发panic
	case TsmsConfig:
		smsConf := conf.(TsmsConfig)
		s.TencentSms = &tencentSms{
			config: &smsConf,
		}
		typ = confType
	}
	return s, typ
}
