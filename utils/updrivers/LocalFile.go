// 上传文件到本地目录
package updrivers

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/perpower/goframe/funcs"

	"github.com/gin-gonic/gin"
)

type Localfile struct{}

// 上传文件
func (*Localfile) Upload(c *gin.Context, objectKey string, filePath *multipart.FileHeader) error {
	err := c.SaveUploadedFile(filePath, funcs.GetRootPath()+objectKey)

	if err != nil {
		return errors.New("Put file failed,error info:" + err.Error())
	}
	return nil
}

// 删除单个文件
func (*Localfile) Delete(objectKey string) error {
	if err := os.Remove(objectKey); err != nil {
		return errors.New("本地文件删除失败, err:" + err.Error())
	}

	return nil
}
