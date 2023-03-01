// 腾讯云COS对象存储操作
package updrivers

import (
	"context"
	"errors"
	"go-framework/configs"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/perpower/goframe/funcs"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentCOS struct{}

// 使用高级上传接口上传对象，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  对象键
// file: *multipart.FileHeader 本地文件
func (*TencentCOS) Upload(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient()

	filePath, err1 := funcs.CreateTempPath(file)
	if err1 != nil {
		return "", err1
	}

	uploadResult, _, err := client.Object.Upload(
		context.Background(), objectKey, filePath, &cos.MultiUploadOptions{
			PartSize: configs.UploadConfigs["tencentCos"]["partSize"].(int64),
		},
	)

	if err != nil {
		return "", err
	}

	if uploadResult.Location != "" {
		speedUrl := configs.UploadConfigs["tencentCos"]["speedUrl"].(string)
		if speedUrl != "" {
			return configs.UploadConfigs["tencentCos"]["speedUrl"].(string) + "/" + objectKey, nil
		}
		return configs.UploadConfigs["tencentCos"]["defaultUrl"].(string) + "/" + objectKey, nil
	} else {
		return "", nil
	}
}

// 使用简单上传接口
// objectKey: string  对象键
// file: *multipart.FileHeader 本地文件
func (*TencentCOS) Put(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient()

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err2 := client.Object.Put(context.Background(), objectKey, f, nil)
	if err2 != nil {
		return "", errors.New("Put file failed,error info:" + err2.Error())
	}

	speedUrl := configs.UploadConfigs["tencentCos"]["speedUrl"].(string)
	if speedUrl != "" {
		return configs.UploadConfigs["tencentCos"]["speedUrl"].(string) + "/" + objectKey, nil
	}
	return configs.UploadConfigs["tencentCos"]["defaultUrl"].(string) + "/" + objectKey, nil

}

// 删除文件对象
func (*TencentCOS) Delete(objectKey string) error {
	client := newClient()

	fullName := configs.UploadConfigs["tencentCos"]["defaultUrl"].(string) + "/" + objectKey
	_, err := client.Object.Delete(context.Background(), fullName)
	if err != nil {
		return errors.New("Delete file failed,error info:" + err.Error())
	}
	return nil
}

// 初始化COS client
func newClient() *cos.Client {
	u, _ := url.Parse(configs.UploadConfigs["tencentCos"]["defaultUrl"].(string))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			SecretID: configs.UploadConfigs["tencentCos"]["secretId"].(string),
			// 环境变量 SECRETKEY 获取用户的 SecretKey
			SecretKey: configs.UploadConfigs["tencentCos"]["secretKey"].(string),
		},
	})

	return client
}
