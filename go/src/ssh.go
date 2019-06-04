package main

import (
	"./utils"
	"fmt"
)

func main() {

	var args map[string]string

	args = utils.GetArgs()

	if args == nil {
		return
	}

	switch args["operation"] {
	case "create":
		utils.SinAddSSH(args)
	case "delete":
		utils.SinDelSSH(args)
	case "connect":
		utils.SinConSSH(args)
	default:
		fmt.Println("Error operation!")
	}
	return
}
