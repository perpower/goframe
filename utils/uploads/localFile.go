// 上传文件到本地目录
package uploads

import (
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

// 本地上传配置结构体
type LocalConfig struct {
	UploadDir string // 文件存储地址
}

type localFile struct{}

// Upload 上传文件
func (*localFile) Upload(c *gin.Context, objectUrl string, file *multipart.FileHeader) error {
	return c.SaveUploadedFile(file, objectUrl)
}

// 删除单个文件
func (*localFile) Delete(objectUrl string) error {
	return os.Remove(objectUrl)
}
