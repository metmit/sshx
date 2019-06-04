package utils

import (
	"fmt"
	"../common"
	"os"
)

//文件名
var FileName string = ""


/**
 * 获取存储文件
 */
func GetStoreFile(file_name string) string {

	if FileName != "" {
		return FileName
	}

	// 用户目录
	home,err := Home()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	path := home + common.FOLDER_NAME

	FileName = path + file_name

	return FileName
}


// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}