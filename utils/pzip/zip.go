// 文件压缩功能组件
package pzip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/perpower/goframe/utils/pfile"
)

var Zip = gzip{}

type gzip struct{}

// 压缩N个指定文件到目标压缩文件中(该方法支持符号链接)
// If the specified files or dirs is a symbolic link ZipFollowSymlink will follow it.
// Note that the symbolic link need to avoid loops.
func (z *gzip) ZipFollowSymlink(zipPath string, paths ...string) error {
	// Create zip file and it's parent dir.
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// New zip writer.
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Get all the file or directory paths.
	var allFilePaths []string
	pathToRoot := make(map[string]string)
	for _, path := range paths {
		// If the path is a dir or symlink to dir, get all files in it.
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Remove the trailing path separator if path is a directory.
			path = strings.TrimSuffix(path, string(os.PathSeparator))
			filePaths, err := pfile.Get.ListDirEntryPathsSymlink(path, true)
			if err != nil {
				return err
			}
			allFilePaths = append(allFilePaths, filePaths...)
			for _, p := range filePaths {
				pathToRoot[p] = path
			}
			continue
		}
		allFilePaths = append(allFilePaths, path)
		pathToRoot[path] = path
	}

	// Traverse all the file or directory.
	for _, path := range allFilePaths {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		// Create a local file header.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set compression method.
		header.Method = zip.Deflate

		// Set relative path of a file as the header name.
		header.Name, err = filepath.Rel(filepath.Dir(pathToRoot[path]), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		// Create writer for the file header and save content of the file.
		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// If file is a directory.
		if info.IsDir() {
			continue
		}

		// If file is a file or symlink to file.
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		f, err := os.Open(realPath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// 压缩N个指定文件到目标压缩文件中(该方法不支持符号链接)
// If a path is a dir don't need to specify the trailing path separator.
// For example calling Zip("archive.zip", "dir", "csv/baz.csv") will get archive.zip and the content of which is
// baz.csv
// dir
// ├── bar.txt
// └── foo.txt
// Note that if a file is a symbolic link it will be skipped.
// zipPath: string 目标压缩文件路径
func (z *gzip) Zip(zipPath string, paths ...string) error {
	// Create zip file and it's parent dir.
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// New zip writer.
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Traverse the file or directory.
	for _, rootPath := range paths {
		// Remove the trailing path separator if path is a directory.
		rootPath = strings.TrimSuffix(rootPath, string(os.PathSeparator))

		// Visit all the files or directories in the tree.
		err = filepath.Walk(rootPath, z.walkFunc(rootPath, zipWriter))
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *gzip) walkFunc(rootPath string, zipWriter *zip.Writer) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If a file is a symbolic link it will be skipped.
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Create a local file header.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set compression method.
		header.Method = zip.Deflate

		// Set relative path of a file as the header name.
		header.Name, err = filepath.Rel(filepath.Dir(rootPath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		// Create writer for the file header and save content of the file.
		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(headerWriter, f)
		return err
	}
}
