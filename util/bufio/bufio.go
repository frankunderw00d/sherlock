package bufio

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type ()

const ()

var ()

func init() {}

// 按指定分割符分割指定文件
// path : 指定文件路径
// sep : 指定分割符
func SplitFile(path string, sep byte) ([]string, error) {
	if path == "" {
		return nil, errors.New("path can't be nil")
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	l := make([]string, 0)
	for {
		str, err := reader.ReadString(sep)
		if err != nil {
			if len(str) > 0 {
				l = append(l, str)
			}
			break
		}
		l = append(l, str)
	}

	return l, file.Close()
}

// 指定从 reader 中逐行读取
func Scan(reader io.Reader) []string {
	scanner := bufio.NewScanner(reader)

	l := make([]string, 0)
	for scanner.Scan() {
		l = append(l, scanner.Text())
	}
	return l
}

// 带缓存的写入，减少IO操作
// path : 指定文件路径
// str... : 写入数据
func Write(path string, data ...string) (int64, error) {
	if path == "" || len(data) <= 0 {
		return 0, errors.New("path or data can't be nil")
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return 0, err
	}

	writer := bufio.NewWriter(file)

	totalLen := int64(0)
	var e error
	for _, str := range data {
		n, err := writer.WriteString(str)
		if n >= 0 {
			totalLen += int64(n)
		}
		if err != nil {
			e = err
			break
		}
	}
	if e != nil {
		return totalLen, e
	}

	if err := writer.Flush(); err != nil {
		return totalLen, err
	}

	return totalLen, file.Close()
}
