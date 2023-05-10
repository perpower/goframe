// 腾讯云COS对象存储操作
package uploads

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

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

type Object = cos.Object
type ObjectTag = cos.ObjectTaggingTag

type tencentCos struct {
	config *CosConfig
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

// 使用高级上传接口上传对象，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  文件对象
// file: *multipart.FileHeader 本地文件
func (t *tencentCos) Upload(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient(t.config)

	filePath, err := dpath.CreateTempPath(file)
	if err != nil {
		return "", err
	}

	uploadResult, _, err := client.Object.Upload(
		context.TODO(), objectKey, filePath, &cos.MultiUploadOptions{
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
// objectKey: string  文件对象
// file: *multipart.FileHeader 本地文件
func (t *tencentCos) Put(objectKey string, file *multipart.FileHeader) (string, error) {
	client := newClient(t.config)

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err1 := client.Object.Put(context.TODO(), objectKey, f, nil)
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
// objectKey: string 文件对象
func (t *tencentCos) Delete(objectKey string) error {
	client := newClient(t.config)

	fullName := t.config.DefaultUrl + "/" + objectKey
	_, err := client.Object.Delete(context.TODO(), fullName)
	return err
}

// 下载对象到本地目录
// objectKey: string 文件对象
// filepath: string 本地文件路径
func (t *tencentCos) Download(objectKey, filepath string) error {
	client := newClient(t.config)

	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}
	_, err := client.Object.Download(context.TODO(), objectKey, filepath, opt)
	return err
}

// 下载对象到浏览器弹窗下载
// objectKey: string 文件对象
// fileName: string 下载保存的文件名
func (t *tencentCos) DownloadWeb(w http.ResponseWriter, objectKey, fileName string) {
	client := newClient(t.config)

	resp, err := client.Object.Get(context.TODO(), objectKey, nil)
	if err != nil {
		http.Error(w, "服务器异常", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "禁止的请求", http.StatusForbidden)
		return
	}

	// 流式下载
	t.DownloadStreamFile(w, resp.Body, fileName)
}

// 初始化分片上传, 并返回uploadID
// objectKey: string  文件对象
func (t *tencentCos) InitiateMultipartUpload(objectKey string) (string, error) {
	client := newClient(t.config)
	v, _, err := client.Object.InitiateMultipartUpload(context.TODO(), objectKey, nil)
	if err != nil {
		return "", err
	}
	return v.UploadID, nil
}

// 分片上传(注意：块有最小1MB的限制，小文件不要用分块上传)
// objectKey: string  文件对象
// uploadID: string
// partNumber: int 分片号
// data: []byte 文件数据
// return: partETag, error
func (t *tencentCos) UploadPart(objectKey, uploadID string, partNumber int, data []byte) (string, error) {
	client := newClient(t.config)

	// 转换数据为bytes.Reader
	byteReader := bytes.NewReader(data)
	resp, err := client.Object.UploadPart(context.TODO(), objectKey, uploadID, partNumber, byteReader, &cos.ObjectUploadPartOptions{
		ContentLength: int64(len(data)),
	})

	partETag := resp.Header.Get("ETag")
	return partETag, err
}

// 完成分片上传
// objectKey: string  文件对象
// uploadID: string
// objectParts: []Object
// objectTags: []ObjectTag 对象标签
func (t *tencentCos) CompleteMultipartUpload(objectKey, uploadID string, objectParts []Object, objectTags ...ObjectTag) error {
	client := newClient(t.config)

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = objectParts
	if len(objectTags) > 0 {
		tagStr := ""
		for key, val := range objectTags {
			if key > 0 {
				tagStr += ("&" + val.Key + "=" + val.Value)
			} else {
				tagStr += (val.Key + "=" + val.Value)
			}
		}
		opt.XOptionHeader.Add("x-cos-tagging", tagStr)
	}
	_, _, err := client.Object.CompleteMultipartUpload(context.TODO(), objectKey, uploadID, opt)
	return err
}

// 终止分片上传并删除已上传的块
// objectKey: string  文件对象
// uploadID: string
func (t *tencentCos) AbortMultipartUpload(objectKey, uploadID string) error {
	client := newClient(t.config)

	_, err := client.Object.AbortMultipartUpload(context.TODO(), objectKey, uploadID)
	return err
}

// 判断指定对象是否存在
// objectKey: string  文件对象
// return: isExist, err
func (t *tencentCos) IsExist(objectKey string) (bool, error) {
	client := newClient(t.config)

	ok, err := client.Object.IsExist(context.TODO(), objectKey)
	return ok, err
}

// 给对象设置标签
// objectKey: string  文件对象
// objectTags: []ObjectTag 对象标签
func (t *tencentCos) PutTagging(objectKey string, objectTags ...ObjectTag) error {
	client := newClient(t.config)
	opt := &cos.ObjectPutTaggingOptions{
		TagSet: objectTags,
	}
	_, err := client.Object.PutTagging(context.TODO(), objectKey, opt)
	return err
}

// 删除对象标签
// objectKey: string  文件对象
func (t *tencentCos) DeleteTagging(objectKey string) error {
	client := newClient(t.config)
	_, err := client.Object.DeleteTagging(context.TODO(), objectKey)
	return err
}

// 查询对象标签
// objectKey: string  文件对象
func (t *tencentCos) GetTagging(objectKey string) ([]ObjectTag, error) {
	client := newClient(t.config)
	resp, _, err := client.Object.GetTagging(context.TODO(), objectKey)
	return resp.TagSet, err
}

// 获取预签名URL
// objectKey: string  文件对象
// expired: time.Duration URL有效期
// return: presignedUrl, err
func (t *tencentCos) GetPresignedURL(objectkey string, expired time.Duration) (string, error) {
	client := newClient(t.config)

	presignedURL, err := client.Object.GetPresignedURL(context.TODO(), http.MethodGet, objectkey, t.config.SecretId, t.config.SecretKey, expired, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// 根据预签名URL下载对象
// presignedUrl: string 预签名URL
// fileName: string 下载保存文件名
func (t *tencentCos) DownloadPresignedObject(w http.ResponseWriter, presignedUrl, fileName string) {
	resp, err := http.Get(presignedUrl)
	if err != nil {
		http.Error(w, "服务器异常", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "禁止的请求", http.StatusForbidden)
		return
	}
	// 流式下载
	t.DownloadStreamFile(w, resp.Body, fileName)
}

// 流式下载文件
func (t *tencentCos) DownloadStreamFile(w http.ResponseWriter, respBody io.ReadCloser, fileName string) {
	w.Header().Set("Content-Type", "application/octet-stream") // 默认让浏览器下载文件
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Cache-Control", "no-cache")

	// 使用chunked分块编码
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	// 此处调用Flush()确保第一个数据块能够直接发送给客户端
	w.(http.Flusher).Flush()

	buf := make([]byte, 1024*1024*5) // 设置5MB的缓冲区
	for {
		n, err := respBody.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, "服务器异常", http.StatusInternalServerError)
			return
		}

		if n == 0 {
			// 文件已经读取完毕
			break
		}

		// 将块写入响应体
		if _, err := w.Write(buf[:n]); err != nil {
			http.Error(w, "服务器异常", http.StatusInternalServerError)
			return
		}
		// 刷新缓冲区，确保数据被立即发送
		w.(http.Flusher).Flush()
	}
}
