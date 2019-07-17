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

	con := utils.Connect{Args: args}

	switch args["operation"] {
	case "create":
		con.Add()
	case "delete":
		con.Del()
	case "connect":
		con.Login()
	default:
		fmt.Println("Error operation!")
	}
	return
}
