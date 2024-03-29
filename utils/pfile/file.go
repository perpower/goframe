package pfile

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

var File = gfile{}

type gfile struct{}

// ReadLines reads all lines of the file.
func (pf *gfile) ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadLinesV2 reads all lines of the file.
func (pf *gfile) ReadLinesV2(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	r := bufio.NewReader(file)
	for {
		// ReadString reads until the first occurrence of delim in the input,
		// returning a string containing the data up to and including the delimiter.
		line, err := r.ReadString('\n')
		if err == io.EOF {
			lines = append(lines, line)
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, line[:len(line)-1])
	}
	return lines, nil
}

// ReadLinesV3 reads all lines of the file.
func (pf *gfile) ReadLinesV3(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	r := bufio.NewReader(f)
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, string(bytes))
	}
	return lines, nil
}

// Create creates or truncates the target file specified by path.
// If the file already exists, it is truncated.
// If the parent directory does not exist, it will be created with mode os.ModePerm(0777).
// If the file does not exist, it is created with mode 0666.
// If successfully call Create, the returned file can be used for I/O
// and the associated file descriptor has mode O_RDWR.
func (pf *gfile) Create(path string) (*os.File, error) {
	exist, err := Check.IsExist(path)
	if err != nil {
		return nil, err
	}
	if exist {
		return os.Create(path)
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// CreateFile creates a file specified by path.
// If the file already exists, it is truncated.
// If the parent directory does not exist, it will be created with mode os.ModePerm(0777).
func (pf *gfile) CreateFile(path string) error {
	file, err := pf.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// FileToBytes serialize the file to bytes.
func (pf *gfile) FileToBytes(path string) []byte {
	bytes, _ := os.ReadFile(path)
	return bytes
}

// BytesToFile writes data to a file.
// If the file does not exist it will be created with permission mode 0644.
func (pf *gfile) BytesToFile(path string, data []byte) error {
	exist, _ := Check.IsExist(path)
	if !exist {
		if err := pf.CreateFile(path); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

// ClearFile clears a file content.
func (pf *gfile) ClearFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
