// 文件解压缩功能组件
package pzip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/perpower/goframe/funcs/normal"
)

var Unzip = gunzip{}

type gunzip struct{}

// 文件信息
type FileInfo struct {
	FileName string     // 文件名
	Ext      string     // 扩展名
	Size     int64      // 文件大小
	Isdir    bool       // 是否文件夹
	Children []FileInfo // 子级内容
}

// 读取指定压缩包内的文件
// filepath: string 文件路径
func (z *gunzip) ReadContent(filename string) ([]FileInfo, error) {
	// 打开一个zip格式文件
	r, err := zip.OpenReader(filename)
	if err != nil {
		return []FileInfo{}, err
	}
	defer r.Close()

	if len(r.File) <= 0 {
		return []FileInfo{}, nil
	}

	// 迭代压缩文件中的文件，打印出文件信息
	fileList := make([]FileInfo, len(r.File))
	for k, f := range r.File {
		fileList[k] = FileInfo{
			FileName: normal.ConvertToUtf8(f.Name),
			Ext:      getFileExt(f.Name),
			Size:     f.FileInfo().Size(),
			Isdir:    f.FileInfo().IsDir(),
		}
	}

	return fileList, nil
}

// 判断是否是目录
func (z *gunzip) IsDir(content []byte) bool {
	info, err := os.Stat(string(content))
	if os.IsNotExist(err) {
		// path does not exist
		return false
	}
	return info.IsDir()
}

// 判断文件是否是zip文件
func (z *gunzip) IsZip(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".zip" || ext == ".ZIP"
}

// 获取文件扩展名
// fileName: string 文件名称
func getFileExt(fileName string) (ext string) {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		ext = parts[len(parts)-1]
	}
	return
}

// 解压缩zip文件到指定目录
// 注意：目标目录后面不要带“/”
// filepath: string zip文件路径
// dstDir: string 解压目标目录
func (z *gunzip) Unzip(filepath, dstDir string) error {
	// open zip file
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := z.unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

func (z *gunzip) unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
