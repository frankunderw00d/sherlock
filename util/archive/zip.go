// zip 包提供了 zip 档案文件的读写服务
package archive

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 压缩文件
// path : 压缩到某个路径下的文件名
// files : 指定的文件集
func ZipCompress(path string, files ...string) error {
	if path == "" || files == nil || len(files) == 0 {
		return errors.New("files can't be nil or empty")
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	zw := zip.NewWriter(file)

	for _, subFilePath := range files {
		subFile, err := os.OpenFile(subFilePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		subFileInfo, err := subFile.Stat()
		if err != nil {
			return err
		}

		fh, err := zip.FileInfoHeader(subFileInfo)
		if err != nil {
			return err
		}
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, subFile)
		if err != nil {
			return err
		}
		subFile.Close()
	}

	return zw.Close()
}

// 解压缩文件
// path : 压缩文件路径
// dirPath : 解压缩到指定文件夹下
func ZipDecompress(path string, dirPath string) error {
	if path == "" || dirPath == "" {
		return errors.New("path or dirPath can't be empty")
	}

	zw, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zw.Close()

	for _, file := range zw.File {
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(file.Name, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		fr, err := file.Open()
		if err != nil {
			return err
		}

		fw, err := os.OpenFile(filepath.Join(dirPath, file.Name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		fw.Close()
		fr.Close()
	}
	return nil
}
