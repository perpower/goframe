// 邮件发送类
package pmailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"reflect"

	"github.com/perpower/goframe/funcs/dpath"

	"gopkg.in/gomail.v2"
)

// 发件服务传参结构体
type EmailSererConfig struct {
	ServerAddress string //邮件服务地址
	Port          int    //邮件端口
	Username      string //发件人地址
	Password      string //密码
}

// 附件传参结构体
type AttachFormat struct {
	Filename string
	Settings []gomail.FileSetting
}

// 收件传参结构体
type EmailConfig struct {
	To      []string       //收件人地址
	Cc      []string       //抄送地址
	Subject string         //邮件主题
	Body    string         //邮件内容,html格式
	Attach  []AttachFormat //附件
}

// newDialer 实例化SMTP Dialer
// smtp: 发件服务器配置
func newDialer(smtp EmailSererConfig) *gomail.Dialer {
	d := gomail.NewDialer(smtp.ServerAddress, smtp.Port, smtp.Username, smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}

// Send 发送邮件
// smtp: 发件服务器配置
// params: 收件人配置
func Send(smtp EmailSererConfig, params EmailConfig) (bool, error) {
	d := newDialer(smtp)

	checkConfig := reflect.TypeOf(params)
	m := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	m.SetHeader("From", smtp.Username)
	m.SetHeader("To", params.To...)
	if _, exist := checkConfig.FieldByName("Cc"); exist && len(params.Cc) > 0 {
		m.SetHeader("Cc", params.Cc...)
	}
	m.SetHeader("Subject", params.Subject)
	m.SetBody("text/html", params.Body)
	if _, exist := checkConfig.FieldByName("Attach"); exist && len(params.Attach) > 0 {
		for _, file := range params.Attach {
			path := dpath.IsPathExist(file.Filename)
			if !path {
				fmt.Println("Error:", file.Filename, "does not exist")
				continue
			} else {
				fmt.Println("uploading", file.Filename, "...")
				m.Attach(file.Filename, file.Settings...)
			}
		}
	}

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}

// GetTplContentByFile 获取指定邮件模板内容
// tplPath: string 邮件模板文件地址
// mailData: interface{}  解析后的邮件内容
func GetTplContentByFile(tplPath string, mailData interface{}) (string, error) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// GetEmailHTMLContent 获取邮件模板
// mailTpl: string 邮件模板路径
// mailData: interface{} 邮件默认需要传递的数据
func GetEmailHTMLContent(mailTpl string, mailData interface{}) (string, error) {
	tpl, err := template.New("emailTpl.html").Parse(mailTpl)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
