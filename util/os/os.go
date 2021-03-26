package os

import "os"

// 存在 且 文件
func FileExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 存在 且 文件夹
func DirExist(name string) bool {
	info, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
