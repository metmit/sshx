package utils

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Expect struct {
}

func (e *Expect) Connect(addr string, port string, user string, pass string) {

	//设置环境变量
	_ = os.Setenv("SIN_TERM_HOST", addr)
	_ = os.Setenv("SIN_TERM_PORT", port)
	_ = os.Setenv("SIN_TERM_USER", user)
	_ = os.Setenv("SIN_TERM_PASS", pass)

	//生成临时文件
	tempFile := e.createExp()

	//执行
	params := e.cmdPath("expect") + ` -f ` + tempFile + ` && ` + e.cmdPath("rm") + ` -rf ` + tempFile

	// exit 时无法删除
	defer os.Remove(tempFile)

	sh := e.cmdPath("sh")
	args := []string{sh, "-c", params}
	env := os.Environ()
	_ = syscall.Exec(sh, args, env)
}

func (e *Expect) createExp() string {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "sinTerm-*.exp")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	content := `#!` + e.cmdPath("expect") + ` -f  
set timeout -1
send_user " \n"
spawn ` + e.cmdPath("ssh") + ` -o StrictHostKeyChecking=no -p $env(SIN_TERM_PORT) $env(SIN_TERM_USER)@$env(SIN_TERM_HOST);
expect {
"*yes/no" { send "yes\r"; exp_continue}
"*password:" { send "$env(SIN_TERM_PASS)\r" }
}
interact`

	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return tmpFile.Name()
}

func (e *Expect) cmdPath(cmd string) string {
	binary, lookErr := exec.LookPath(cmd)
	if lookErr != nil {
		panic(lookErr)
	}
	return binary
}
