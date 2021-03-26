// tar 包实现了 tar 格式压缩文件的存取
package archive

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 压缩文件
// path : 压缩到某个路径下的文件名
// files : 指定的文件集
func TarCompress(path string, files ...string) error {
	if path == "" || files == nil || len(files) == 0 {
		return errors.New("files can't be nil or empty")
	}

	// 创建写入文件
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建压缩写入
	tw := tar.NewWriter(file)

	// 遍历 files 路径
	for _, filePath := range files {
		// 读取
		subFile, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		// 读取信息
		subFileInfo, err := subFile.Stat()
		if err != nil {
			return err
		}

		// 写入头部
		err = tw.WriteHeader(&tar.Header{
			Name:    subFileInfo.Name(),
			Size:    subFileInfo.Size(),
			Mode:    int64(subFileInfo.Mode()),
			ModTime: subFileInfo.ModTime(),
		})
		if err != nil {
			return err
		}

		// 读取全部数据
		data, err := ioutil.ReadAll(subFile)
		if err != nil {
			return err
		}

		// 写入数据
		_, err = tw.Write(data)
		if err != nil {
			return err
		}

		// 关闭文件
		err = subFile.Close()
		if err != nil {
			return err
		}
	}

	// 关闭写入
	return tw.Close()
}

// 解压缩文件
// path : 压缩文件路径
// dirPath : 解压缩到指定文件夹下
func TarDecompress(path string, dirPath string) error {
	if path == "" || dirPath == "" {
		return errors.New("path or dirPath can't be empty")
	}

	// 打开读取文件
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// 新建读取
	tr := tar.NewReader(file)

	for {
		hdr, err := tr.Next() // 切换下一个
		if err == io.EOF {    // EOF 结束
			break
		}
		if err != nil {
			return err
		}

		// 读取头部的信息，创建和写入文件信息
		subFile, err := os.OpenFile(filepath.Join(dirPath, hdr.Name), os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
		if err != nil {
			return err
		}
		// 写入信息
		_, err = io.Copy(subFile, tr)
		if err != nil {
			return err
		}
		// 关闭
		err = subFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
