// 文件上传能力
package utils

import (
	"go-framework/configs"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/perpower/goframe/funcs"
	"github.com/perpower/goframe/utils/updrivers"

	"github.com/gin-gonic/gin"
)

// 统一上传方法
// uploadType: string 指定上传文件的方式
// fileType: string 上传文件类型
func Upload(c *gin.Context, uploadType string, fileType string) (string, error) {
	if _, ok := configs.UploadConfigs[uploadType]; !ok {
		log.Println("上传方式配置有误")
		return "", nil
	}

	file, err := c.FormFile("file")

	//判断上传文件是否为空
	if file == nil {
		log.Println("请选择文件, errer info:", err)
		return "", err
	}

	_, exist := configs.AllowExt[fileType]
	if exist {
		fileExt := path.Ext(file.Filename)
		if !funcs.InArray(fileExt, configs.AllowExt[fileType]) {
			log.Printf("文件格式%s不在允许上传的范围内", fileExt)
			return "", nil
			//Todo 返回统一的错误处理
		}
	}
	if configs.MaxLimitSize > 0 {
		if file.Size > configs.MaxLimitSize {
			log.Println("文件大小超出允许上传限制")
			return "", nil
			//Todo 返回统一的错误处理
		}
	}

	year, month, day := time.Now().Date()

	switch uploadType {
	case "localFile":
		//判断文件目录是否存在，不存在则创建
		dir := path.Join(configs.UploadConfigs[uploadType]["uploadDir"].(string), strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
		err := funcs.Mkmultdir(dir)
		if err == nil {
			objectUrl := path.Join(dir, file.Filename)
			localFile := updrivers.Localfile{}
			res := localFile.Upload(c, objectUrl, file)
			if res == nil {
				log.Println("文件上传成功")
				return objectUrl, nil
			}

			log.Println("文件上传失败,error info：", res)
			return "", res
		}

		return "", err
	case "tencentCos":
		objectKey := path.Join(configs.UploadConfigs[uploadType]["folder"].(string), strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day), file.Filename)
		tencentCos := updrivers.TencentCOS{}

		objectUrl, err3 := tencentCos.Upload(objectKey, file)

		if objectUrl != "" {
			log.Println("文件上传成功")
			return objectUrl, nil
		}

		log.Println("文件上传失败,error info：", err3)
		return "", err3

	default:
		log.Println("文件上传失败,上传方式错误")
		return "", nil
	}

}
