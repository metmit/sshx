package utils

import (
	"flag"
	"fmt"
)

const SSH_NAME string = "ssh"

var (
	help      bool
	version   bool
	operation string
	cname     string
	secret    string
)

func GetArgs() map[string]string {

	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&version, "v", false, "show version")

	flag.StringVar(&operation, "o", "connect", "operation:[add|create|del|delete|con|connect], default:con")

	flag.StringVar(&cname, "n", "", "Connect name")

	flag.StringVar(&secret, "s", "", "Connect secret")

	// 改变默认的 Usage
	flag.Usage = usage

	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	if version {
		fmt.Println("1.0.0")
		return nil
	}

	cname = "sin"
	for {
		if cname != "" {
			break;
		}
		fmt.Println("Type Connect Name: ")
		_, _ = fmt.Scanln(&cname)
	}

	secret = "sin"

	for {
		if secret != "" {
			break;
		}
		fmt.Println("Type Connect Secret: ")
		_, _ = fmt.Scanln(&secret)
	}

	if operation == "del" {
		operation = "delete"
	}

	if operation == "con" {
		operation = "connect"
	}

	if operation == "add" {
		operation = "create"
	}

	result := map[string]string{
		"name":      cname,
		"secret":    secret,
		"operation": operation,
	}

	return result
}

func usage() {
	fmt.Println("NAME")
	fmt.Println(SSH_NAME + " -- safe ssh soft.")
	fmt.Println("")
	fmt.Println("SYNOPSIS")
	fmt.Println(SSH_NAME + " [-n NAME] [-s SECRET] [-o OPERATION]")
	fmt.Println("")
	fmt.Println("DESCRIPTION")
	fmt.Println(SSH_NAME + " can save ssh host,port,user,password to local file by name & secret, quick resolve and connect to server when 'connect' operation.")
	fmt.Println("Because of resolve depends on the name & secret, so it's safe!")
	fmt.Println("")
	fmt.Println("The options are as follows:")
	flag.PrintDefaults()
}
