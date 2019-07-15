package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
		fmt.Println("文件已存在")
		fmt.Println(c.Args["name"] + "的文件已存在，是否覆盖[Y]: ")
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
			break;
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
			break;
		}
	}

	//# 加密账号密码
	host = base64.StdEncoding.EncodeToString([]byte(host))
	port = base64.StdEncoding.EncodeToString([]byte(port))
	user = base64.StdEncoding.EncodeToString([]byte(user))
	pass = base64.StdEncoding.EncodeToString([]byte(pass))
	// 拼接
	params := host + "@" + port + "@" + user + "@" + pass + "@"

	result := s.Encode(params, c.Args["secret"])
	if result == "" {
		fmt.Println("Add fail!")
		return
	}

	config := Config{}
	result = "v" + config.Version + "v" + result

	fmt.Println(result)
	fmt.Println(fileName)

	if ioutil.WriteFile(fileName, []byte(result), 0644) != nil {
		fmt.Println("Write file fail!")
	}
	return
}

//删除链接
func (c *Connect) Del() {

}

//连接
func (c *Connect) Login() {

	f := GetFileInstance();
	s := GetStrInstance();
	fileName := f.GetFullName(c.Args["name"])

	// 如果文件不存在
	if !f.Exists(fileName) {
		fmt.Println("文件不存在")
		return
	}

	byteContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("读取文件失败")
		return
	}

	content := strings.Replace(string(byteContent), "\n", "", 1)
	content = strings.Replace(content, "\\ ", "", 1)
	content = strings.Replace(content, " ", "", 1)
	content = content[strings.LastIndex(content, "v")+1:]

	base := s.Decode(content, c.Args["secret"])

	split := strings.Split(base, "@")

	if len(split) < 4 {
		fmt.Println("Has Not Password!")
		return
	}

	host, _ := base64.StdEncoding.DecodeString(split[0])
	port, _ := base64.StdEncoding.DecodeString(split[1])
	user, _ := base64.StdEncoding.DecodeString(split[2])
	pass, _ := base64.StdEncoding.DecodeString(split[3])

	sHost := strings.Replace(string(host), "\n", "", 1)
	sPort := strings.Replace(string(port), "\n", "", 1)
	sUser := strings.Replace(string(user), "\n", "", 1)
	sPass := strings.Replace(string(pass), "\n", "", 1)

	expect := Expect{}
	expect.Connect(sHost, sPort, sUser, sPass)
}
