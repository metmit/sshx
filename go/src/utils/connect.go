package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Connect struct {
	Args map[string]string
}

//添加链接
func (c *Connect) Add() {

	f := GetFileInstance()
	s := GetStrInstance()
	fileName := f.GetFullName(c.Args["name"])

	// 如果文件已存在
	if f.Exists(fileName) {
		var rewrite string
		fmt.Println(c.Args["name"] + " file already exist, rewrite?[Y]: ")
		_, _ = fmt.Scanln(&rewrite)
		if rewrite == "" {
			rewrite = "Y"
		}
		if strings.ToUpper(rewrite) != "Y" {
			return
		}
	}

	var (
		host string
		port string
		user string
		pass string
	)

	for {
		fmt.Println("host: ")
		_, _ = fmt.Scanln(&host)
		if host != "" {
			break
		}
	}

	fmt.Println("port[22]: ")
	_, _ = fmt.Scanln(&port)
	if port == "" {
		port = "22"
	}

	fmt.Println("user[root]: ")
	_, _ = fmt.Scanln(&user)
	if user == "" {
		user = "root"
	}

	for {
		fmt.Println("password: ")
		_, _ = fmt.Scanln(&pass)
		if pass != "" {
			break
		}
	}

	// 拼接
	params := s.Base64Encode(host) + "@" + s.Base64Encode(port) + "@" + s.Base64Encode(user) + "@" + s.Base64Encode(pass) + "@"

	result := s.Encode(params, c.Args["secret"])
	if result == "" {
		fmt.Println("Encode content fail!")
		return
	}

	result = "v" + GetConfig().Version + "v" + result

	if ioutil.WriteFile(fileName, []byte(result), 0644) != nil {
		fmt.Println("write file fail!")
	}
	return
}

//删除链接
func (c *Connect) Del() {
	f := GetFileInstance()

	fileName := f.GetFullName(c.Args["name"])

	if !f.Exists(fileName) {
		fmt.Println("File dos not exist, delete fail!")
		return
	}

	var rewrite string
	fmt.Println(c.Args["name"] + " delete ? [N|Y]: ")
	_, _ = fmt.Scanln(&rewrite)
	if rewrite == "" {
		rewrite = "N"
	}
	if strings.ToUpper(rewrite) == "Y" {
		_ = os.Remove(fileName)
		return
	}
}

//连接
func (c *Connect) Login() {

	f := GetFileInstance()
	s := GetStrInstance()
	fileName := f.GetFullName(c.Args["name"])

	// 如果文件不存在
	if !f.Exists(fileName) {
		fmt.Println("File dos not exist, please add!")
		return
	}

	byteContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Get file content fail!")
		return
	}

	content := strings.Replace(string(byteContent), "\n", "", 20)
	content = strings.Replace(content, "\\ ", "", 20)
	content = strings.Replace(content, " ", "", 20)
	content = content[strings.LastIndex(content, "v")+1:]

	base := s.Decode(content, c.Args["secret"])

	split := strings.Split(base, "@")

	if len(split) < 4 {
		fmt.Println("Has Not Password!")
		return
	}

	expect := Expect{}
	expect.Connect(s.Base64Decode(split[0]), s.Base64Decode(split[1]), s.Base64Decode(split[2]), s.Base64Decode(split[3]))
}
