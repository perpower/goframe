package dpath

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 获取项目根目录路径
func GetRootPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 创建多级目录
func Mkmultdir(dirpath string) error {
	if !IsPathExist(dirpath) {
		err := os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsPathExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 生成文件本地临时路径
func CreateTempPath(file *multipart.FileHeader) (string, error) {
	//go不提供直接获取上传文件的tempfile,所以只能自己创建临时文件
	tmpfile, _ := file.Open()

	defer tmpfile.Close() //确保程序执行完之后可以被删除

	file2, err := os.CreateTemp("", "temp-")
	if err != nil {
		return "", err
	}
	defer file2.Close() //确保程序执行完之后可以被删除

	_, err2 := io.Copy(file2, tmpfile)
	if cerr := file2.Close(); err2 == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(file2.Name())
		return "", err
	}
	filepath := file2.Name()

	return filepath, nil
}
