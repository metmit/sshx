package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

type File struct {
	FileName string
}

func (f *File) GetFileName(cName string) string {
	return GetStrInstance().Md5(cName) + ".sx"
}

/**
 * 获取存储文件
 */
func (f *File) GetFullName(cName string) string {

	if f.FileName != "" {
		return f.FileName
	}

	// 用户目录
	home, err := f.Home()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	path := home + "/" + GetConfig().FolderName + "/" + GetConfig().Version + "/"
	_ = os.MkdirAll(path, 0644)

	f.FileName = path + f.GetFileName(cName)

	return f.FileName
}

// 判断所给路径文件/文件夹是否存在
func (f *File) Exists(path string) bool {
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
func (f *File) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func (f *File) IsFile(path string) bool {
	return !f.IsDir(path)
}

func (f *File) Home() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}

	if "windows" == runtime.GOOS {
		return f.homeWindows()
	}

	return f.homeUnix()
}

func (f *File) homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func (f *File) homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
