package pfile

import (
	"os"
	"path/filepath"
)

var Check = gcheck{}

type gcheck struct{}

// IsExist checks whether the file or directory exists.
// Use os.Stat to get the info of the target file or dir to check whether exists.
// If os.Stat returns nil err, the target exists.
// If os.Stat returns a os.ErrNotExist err, the target does not exist.
// If the error returned is another type, the target is uncertain whether exists.
func (pf *gcheck) IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir checks whether a path is a directory.
// If the path is a symbolic link will follow it.
func (pf *gcheck) IsDir(path string) bool {
	is, _ := pf.IsDirE(path)
	return is
}

// IsDirE checks whether a path is a directory with error.
// If the path is a symbolic link will follow it.
func (pf *gcheck) IsDirE(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true, nil
	}
	return false, err
}

// IsFile checks whether a path is a file.
// If the path is a symbolic link will follow it.
func (pf *gcheck) IsFile(path string) bool {
	is, _ := pf.IsFileE(path)
	return is
}

// IsFileE checks whether a path is a file with error.
// If the path is a symbolic link will follow it.
func (pf *gcheck) IsFileE(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.Mode().IsRegular() {
		return true, nil
	}
	return false, err
}

// IsSymlink checks a file whether is a symbolic link on Linux.
// Note that this doesn't work for the shortcut file on windows.
// If you want to check a file whether is a shortcut file on Windows please use IsShortcut function.
func (pf *gcheck) IsSymlink(path string) bool {
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}

// IsSymlinkE checks a file whether is a symbolic link on Linux.
// Note that this doesn't work for the shortcut file on windows.
// If you want to check a file whether is a shortcut file on Windows please use IsShortcut function.
func (pf *gcheck) IsSymlinkE(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		return true, nil
	}
	return false, err
}

// IsShortcut checks a file whether is a shortcut on Windows.
func (pf *gcheck) IsShortcut(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".lnk" {
		return true
	}
	return false
}
