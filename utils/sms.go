// 短信能力
package utils

import (
	"go-framework/configs"
	"log"

	"github.com/perpower/goframe/utils/smsdrivers"
)

func SendSms(smsType string, templateInfo map[string]interface{}) (bool, error) {
	if _, ok := configs.SmsConfigs[smsType]; !ok {
		log.Println("短信平台配置有误")
		return false, nil
	}

	switch smsType {
	case "tencentSms":
		smsconfig := smsdrivers.SmsConfig{
			SecretId:         configs.SmsConfigs[smsType]["secretId"].(string),
			SecretKey:        configs.SmsConfigs[smsType]["secretKey"].(string),
			AppId:            configs.SmsConfigs[smsType]["appId"].(string),
			SignName:         configs.SmsConfigs[smsType]["signName"].(string),
			TemplateId:       templateInfo["templateId"].(string),
			TemplateParamSet: templateInfo["paramArr"].([]string),
			PhoneNumberSet:   templateInfo["phoneArr"].([]string),
		}

		if _, ok := templateInfo["extraInfo"]; ok {
			smsconfig.SessionContext = templateInfo["extraInfo"].(string)
		}
		if _, ok := templateInfo["extendCode"]; ok {
			smsconfig.SessionContext = templateInfo["extendCode"].(string)
		}
		if _, ok := templateInfo["senderId"]; ok {
			smsconfig.SessionContext = templateInfo["senderId"].(string)
		}

		tencentSms := &smsdrivers.TencentSms{}
		status, _, err := tencentSms.SendSms(smsconfig)

		return status, err
	default:
		return false, nil
	}
}
