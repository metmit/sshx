package utils

import (
	"flag"
	"fmt"
)

func GetArgs() map[string]string {

	var (
		help      bool
		version   bool
		operation string
		cname     string
		secret    string
	)

	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&version, "v", false, "show version")

	flag.StringVar(&operation, "o", "con", "operation:[add|del|con]")

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
		fmt.Println(GetConfig().Version)
		return nil
	}

	for {
		if cname != "" {
			break
		}
		fmt.Println("Type Connect Name: ")
		_, _ = fmt.Scanln(&cname)
	}

	for {
		if secret != "" {
			break
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
	fmt.Println(GetConfig().Name + " -- safe ssh soft.")
	fmt.Println("")
	fmt.Println("SYNOPSIS")
	fmt.Println(GetConfig().Name + " [-n NAME] [-s SECRET] [-o OPERATION]")
	fmt.Println("")
	fmt.Println("DESCRIPTION")
	fmt.Println(GetConfig().Name + " save the encrypted content(host,port,user,password) to a local file, easy resolve and connect to server.")
	fmt.Println("Because of resolve depends on the name & secret, so it's safe!")
	fmt.Println("")
	fmt.Println("The options are as follows:")
	flag.PrintDefaults()
}
