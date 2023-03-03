// 腾讯云COS对象存储操作
package uploads

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/perpower/goframe/funcs/dpath"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 腾讯云COS对象存储配置结构体
type CosConfig struct {
	AppId      string // 应用ID
	SecretId   string // 秘钥ID
	SecretKey  string // 秘钥值
	Bucket     string // 存储桶名称
	Region     string // 指定地域
	PartSize   int64  // 分片上传块大小，单位MB
	DefaultUrl string // 默认访问地址
	CdnUrl     string // CDN加速地址
	Folder     string // 指定虚拟目录
}

type tencentCos struct {
	config *CosConfig
}

// 使用高级上传接口上传对象，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  对象键
// file: *multipart.FileHeader 本地文件
func (t *tencentCos) Upload(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient(t.config)

	filePath, err := dpath.CreateTempPath(file)
	if err != nil {
		return "", err
	}

	uploadResult, _, err := client.Object.Upload(
		context.Background(), objectKey, filePath, &cos.MultiUploadOptions{
			PartSize: t.config.PartSize,
		},
	)

	if err != nil {
		return "", err
	}

	if uploadResult.Location != "" {
		cdnUrl := t.config.CdnUrl
		if cdnUrl != "" {
			return t.config.CdnUrl + "/" + objectKey, nil
		}
		return t.config.DefaultUrl + "/" + objectKey, nil
	} else {
		return "", nil
	}
}

// 使用简单上传接口
// objectKey: string  对象键
// file: *multipart.FileHeader 本地文件
func (t *tencentCos) Put(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient(t.config)

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err1 := client.Object.Put(context.Background(), objectKey, f, nil)
	if err != nil {
		return "", err1
	}

	cdnUrl := t.config.CdnUrl
	if cdnUrl != "" {
		return t.config.CdnUrl + "/" + objectKey, nil
	}
	return t.config.DefaultUrl + "/" + objectKey, nil

}

// Delete 删除文件对象
// objectKey: string
func (t *tencentCos) Delete(objectKey string) error {
	client := newClient(t.config)

	fullName := t.config.DefaultUrl + "/" + objectKey
	_, err := client.Object.Delete(context.Background(), fullName)
	return err
}

// 初始化COS client
func newClient(conf *CosConfig) *cos.Client {
	u, _ := url.Parse(conf.DefaultUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			SecretID: conf.SecretId,
			// 环境变量 SECRETKEY 获取用户的 ecretKey
			SecretKey: conf.SecretKey,
		},
	})

	return client
}
