// 腾讯云短信
package smsdrivers

import (
	"encoding/json"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSms struct{}

// 定义传参结构体
type SmsConfig struct {
	SecretId         string
	SecretKey        string
	AppId            string   //短信应用appId
	SignName         string   //短信签名
	TemplateId       string   //短信模板ID
	TemplateParamSet []string //模板参数
	PhoneNumberSet   []string //发送手机号
	SessionContext   string   //用户的 session 内容
	ExtendCode       string   //短信码号扩展号
	SenderId         string   //国际/港澳台短信 SenderId
}

// 发送短信
func (*TencentSms) SendSms(SmsConfig SmsConfig) (bool, *sms.SendSmsResponseParams, error) {
	//实例化一个认证对象
	credential := common.NewCredential(
		SmsConfig.SecretId,
		SmsConfig.SecretKey,
	)

	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5

	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	request := sms.NewSendSmsRequest()

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(SmsConfig.AppId)

	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	// 签名信息可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-sign) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-sign) 的签名管理查看
	request.SignName = common.StringPtr(SmsConfig.SignName)

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	// 模板 ID 可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-template) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-template) 的正文模板管理查看
	request.TemplateId = common.StringPtr(SmsConfig.TemplateId)

	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(SmsConfig.TemplateParamSet)

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(SmsConfig.PhoneNumberSet)

	/* 用户的 session 内容（无需要可忽略）: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	request.SessionContext = common.StringPtr(SmsConfig.SessionContext)

	/* 短信码号扩展号（无需要可忽略）: 默认未开通，如需开通请联系 [腾讯云短信小助手] */
	request.ExtendCode = common.StringPtr(SmsConfig.ExtendCode)

	/* 国际/港澳台短信 SenderId（无需要可忽略）: 国内短信填空，默认未开通，如需开通请联系 [腾讯云短信小助手] */
	request.SenderId = common.StringPtr(SmsConfig.SenderId)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)

	// 处理异常
	if terr, ok := err.(*errors.TencentCloudSDKError); ok {
		code := terr.GetCode()
		log.Printf("An API error has returned, error code: %s, error info: %s", code, err)
		return false, response.Response, err
	}
	// 非SDK异常
	if err != nil {
		return false, response.Response, err
	}

	sendStatus := response.Response.SendStatusSet
	jsonData, _ := json.Marshal(response.Response.SendStatusSet)

	log.Printf("请求结果Json格式数据: %s", jsonData)
	if *sendStatus[0].Code != "Ok" {
		return false, response.Response, nil
	}

	return true, response.Response, nil
}
