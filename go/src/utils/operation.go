package utils

import (
	"../common"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	"regexp"
	"syscall"

	//"google.golang.org/grpc/codes"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"
)

//添加链接
func SinAddSSH(args map[string]string) {

	file_name := GetStoreFile(sinMd5(args["name"]) + ".sin")

	// 如果文件已存在
	if Exists(file_name) {
		var rewrite string
		fmt.Println("文件已存在")
		fmt.Println(args["name"] + "的文件已存在，是否覆盖[Y]: ")
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

	var result string

	result = sinEncode(host, port, user, pass, args)
	if result == "" {
		fmt.Println("Add fail!")
		return
	}

	fmt.Println(result)

	if ioutil.WriteFile(file_name, []byte(result), 0644) != nil {
		fmt.Println("Write file fail!")
	}
	return
}


//删除链接
func SinDelSSH(args map[string]string) {

}

//连接
func SinConSSH(args map[string]string) {

	file_name := GetStoreFile(sinMd5(args["name"]) + ".sin")

	// 如果文件不存在
	if !Exists(file_name) {
		fmt.Println("文件不存在")
		return
	}

	byte_content, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Println("读取文件失败")
		return
	}
	content := strings.Replace(string(byte_content), "\n", "", 1)
	content = strings.Replace(content, "\\ ", "", 1)
	content = strings.Replace(content, " ", "", 1)

	content = content[strings.LastIndex(content, "v")+1:]

	base := sinDecode(content, args)

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


	//connectCommand(sHost, sPort, sUser, sPass)

	//connectSession(sHost, sPort, sUser, sPass) //可以运行，但很糟糕

	connectSysCall(sHost, sPort, sUser, sPass)

	//connect(sHost, sPort, sUser, sPass)
}

func connectCommand(addr string, port string, user string, pass string) {

	var stdout bytes.Buffer
	//stdin := bytes.NewBuffer(nil)

	command := "expect -c \"set timeout -1;spawn ssh -o StrictHostKeyChecking=no -p "+port+" "+user+"@"+addr+";expect {*assword:* {send "+pass+"\\r;}}\ninteract\""
	cmd := exec.Command("/bin/bash", "-c", command)

	//command := "ssh -o StrictHostKeyChecking=no -p "+port+" "+user+"@"+addr
	//cmd := exec.Command("/bin/bash", "-c", command)

	//command := "set timeout -1 && ssh -o StrictHostKeyChecking=no -p "+port+" "+user+"@"+addr + " && expect \"password:\" && send \""+pass+"\\r\" && interact "
	//cmd := exec.Command("/usr/bin/expect", command)

	//result, _ := cmd.Output()
	//fmt.Println(string(result))

	cmd.Stdout = &stdout
	//cmdstring := pass
	//cmd.Stdin = stdin
	//in.WriteString(cmdstring)
	if err := cmd.Run(); err != nil {
		return
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return
	}
}

func connectSession(addr string, port string, user string, pass string) {

	client, _ := ssh.Dial("tcp", addr+":"+port, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	session, _ := client.NewSession()
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	session.RequestPty("xterm", 25, 80, modes)

	session.Shell()

	session.Wait()
}

func connect(addr string, port string, user string, pass string) {

	sshClt, err := ssh.Dial("tcp", addr+":"+port, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		log.Fatalf("ssh.Dial(%q) failed: %v", addr, err)
	}
	defer sshClt.Close()

	timeout := 10 * time.Minute
	fmt.Println(timeout)
	timeout = 10 * time.Second
	e, _, err := expect.SpawnSSH(sshClt, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	fmt.Println()

	pass = pass + "\r"
	pass = pass + "\n"
	fmt.Println(pass)

	var stdin string

	for {
		promptRE := regexp.MustCompile(".*")
		result, _, _ := e.Expect(promptRE, timeout)
		fmt.Printf("%s", result)

		fmt.Scanln(&stdin)

		fmt.Printf("aaa %s bbb", stdin)

		e.Send(stdin + "\r\n")
		stdin = ""
	}

	//aa, error := e.ExpectBatch([]expect.Batcher{
	//	&expect.BCas{[]expect.Caser{
	//		&expect.Case{R: regexp.MustCompile(`Welcome`), S:"pwd",T: expect.OK()},
	//		&expect.Case{R: regexp.MustCompile(`Login: `), S: user,
	//			T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong username")), Rt: 3},
	//		&expect.Case{R: regexp.MustCompile(`Password: `), S: pass, T: expect.Next(), Rt: 1},
	//		&expect.Case{R: regexp.MustCompile(`password`), S: pass,
	//			T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong password")), Rt: 1},
	//	}},
	//}, timeout)
	//fmt.Println(aa)
	//fmt.Println(error)

	fmt.Println(term.Greenf("All done"))
}

func connectSysCall(addr string, port string, user string, pass string) {

	binary, lookErr := exec.LookPath("bash")
	if lookErr != nil {
		panic(lookErr)
	}

	//binary = "/usr/bin/expect"
	params := `/usr/bin/expect -c "
        set timeout -1;
        spawn /usr/bin/ssh -o StrictHostKeyChecking=no -p `+port+` `+user+`@`+addr+`; 
        expect {
            *assword:* { 
                send `+pass+`\r; 
            }
        }
        interact 
    "`

	args := []string{binary, "-c", params}
	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

//加密
func sinEncode(host string, port string, user string, pass string, args map[string]string) string {

	secret_hex := strings.ToUpper(sinMd5(args["secret"]))

	secret := hexToBigInt(secret_hex)

	//# 加密账号密码
	host = base64.StdEncoding.EncodeToString([]byte(host))
	port = base64.StdEncoding.EncodeToString([]byte(port))
	user = base64.StdEncoding.EncodeToString([]byte(user))
	pass = base64.StdEncoding.EncodeToString([]byte(pass))

	// 拼接
	str := host + "@" + port + "@" + user + "@" + pass + "@"

	//字符串转16
	h := strings.ToUpper(hex.EncodeToString([]byte(str)))

	//# 转10进制
	info := hexToBigInt(h)

	//# 密文与密钥相加
	info.Add(info, secret)

	//数字变成字符串
	result := info.String()

	return "v"+ common.VERSION +"v" + result
}

func sinDecode(content string, args map[string]string) string {

	//加密密钥
	secret_hex := strings.ToUpper(sinMd5(args["secret"]))
	secret := hexToBigInt(secret_hex)

	content_int := strToBigInt(content)

	//密文与密钥相减
	//private := new(big.Int)
	private := content_int.Sub(content_int, secret)

	//十进制转字符串
	return strings.Replace(string(private.Bytes()), "\n", "", 1)
}

// 十六进制转十进制
func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	//n, _ = n.SetString(hex[2:], 16)
	n, _ = n.SetString(hex, 16)
	return n
}

// 字符串转十进制
func strToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex, 10)
	return n
}

// md5
func sinMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

