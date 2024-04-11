package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// IsPathExist : 路径是否存在
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// SaveFile : 将数据保存到文件
func SaveFile(path string, value string) {
	fileHandler, err := os.Create(path)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)

	fmt.Fprintln(buf, value)

	buf.Flush()
}

// SaveFileAppend : 将数据保存到文件[追加]
func SaveFileAppend(path string, value string) {
	fileHandler, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)
	fmt.Fprintln(buf, value)
	buf.Flush()
}

// ReadFile : 一次性读取文件 string
func ReadFile(filePth string) string {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return ""
	}
	r, err := io.ReadAll(f)
	return string(r)
}

// ReadFileByte : 一次性读取文件 []byte
func ReadFileByte(filePth string) []byte {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil
	}
	r, err := io.ReadAll(f)
	return r
}

// WalkDir : 递归读取目录下所有文件名,存在 *fileList
func WalkDir(dirpath string, fileList *[]string) {
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			WalkDir(dirpath+"/"+file.Name(), fileList)
			continue
		} else {
			*fileList = append(*fileList, dirpath+"/"+file.Name())
		}
	}
}

// ReadListFile : 按行读文件
func ReadListFile(filename string) []string {
	result := make([]string, 0)
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			result = append(result, line)
		}

		if err == io.EOF {
			break
		}
	}
	return result
}

// WriteListFile : 按行存储文件
func WriteListFile(filename string, value []string) {
	fileHandler, err := os.Create(filename)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)

	for _, v := range value {
		fmt.Fprintln(buf, v)
	}
	buf.Flush()
}

// ListDir : 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}

		files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

// IsDir 判断路径是文件夹返回 true，是文件返回 false
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// IsFile 判断路径是文件返回 true，是文件夹返回 false
func IsFile(path string) bool {
	return !IsDir(path)
}
