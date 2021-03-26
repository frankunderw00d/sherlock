package log

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	uos "sherlock/util/os"
	"time"
)

type (
	FileHook struct {
		Suffix        string
		File          *os.File
		Buffer        *bytes.Buffer
		MaxBufferSize int64
		MaxFileSize   int64
	}
)

const (
	MB                          = 1 << 20
	SuggestBufferSize    int64  = 20 * MB   // 建议缓存大小
	SuggestFileSize      int64  = 200 * MB  // 建议文件大小
	SuggestSuffix        string = "log.txt" // 建议文件后缀
	DefaultLogDirName           = "log"
	DefaultDateDirFormat        = "20060102"
	DefaultFileDirFormat        = "20060102150405"
)

var ()

// 新建默认文件钩子
func NewDefaultFileHook() io.WriteCloser {
	fh, _ := NewFileHook(SuggestSuffix, SuggestBufferSize, SuggestFileSize)
	return fh
}

// 新建文件钩子
func NewFileHook(suffix string, mbs, mfs int64) (io.WriteCloser, error) {
	fh := &FileHook{
		Suffix:        suffix,
		File:          nil,
		Buffer:        bytes.NewBuffer([]byte{}),
		MaxBufferSize: mbs,
		MaxFileSize:   mfs,
	}

	if err := fh.exchangeLogFile(); err != nil {
		return nil, err
	}

	return fh, nil
}

// 构建新文件路径
func newLogFilePath(name string) (string, error) {
	now := time.Now()

	dateDir := now.Format(DefaultDateDirFormat)
	timeFile := now.Format(DefaultFileDirFormat) + name

	dirPath := filepath.Join(DefaultLogDirName, dateDir)
	if !uos.DirExist(dirPath) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	return filepath.Join(dirPath, timeFile), nil

}

// 写入
func (fh *FileHook) Write(b []byte) (int, error) {
	if fh.Buffer == nil {
		return 0, errors.New("buffer is nil")
	}
	if fh.File == nil {
		return 0, errors.New("output file is nil")
	}

	n, err := fh.Buffer.Write(b)
	if err != nil {
		return 0, err
	}

	// 判断缓存大小
	if int64(fh.Buffer.Len()) >= fh.MaxBufferSize {
		if _, err := fh.WriteBufferToFile(); err != nil {
			return 0, err
		}

		// 判断文件大小
		fileInfo, err := fh.File.Stat()
		if err != nil {
			return 0, err
		}

		if fileInfo.Size() >= fh.MaxFileSize {
			if err := fh.exchangeLogFile(); err != nil {
				return 0, err
			}
		}
	}

	return n, nil
}

// 关闭
func (fh *FileHook) Close() error {
	if fh.Buffer != nil && fh.Buffer.Len() > 0 {
		if _, err := fh.WriteBufferToFile(); err != nil {
			return err
		}
	}
	if fh.File != nil {
		if err := fh.File.Close(); err != nil {
			return err
		}
	}

	fh.Buffer = nil
	fh.File = nil

	return nil
}

// 将缓存写入文件
func (fh *FileHook) WriteBufferToFile() (int64, error) {
	if fh.Buffer == nil {
		return 0, errors.New("buffer is nil")
	}
	if fh.File == nil {
		return 0, errors.New("output file is nil")
	}

	return fh.Buffer.WriteTo(fh.File)
}

// 关闭旧文件，创建新文件
func (fh *FileHook) exchangeLogFile() error {
	logFilePath, err := newLogFilePath(fh.Suffix)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	if fh.File != nil {
		if err := fh.File.Close(); err != nil {
			return err
		}
	}

	fh.File = file

	return nil
}
