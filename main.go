package main

import (
	"fmt"
	"os"
	"simple-file-server/command"
	"simple-file-server/lib/config"
	"simple-file-server/lib/log"
)

//func isRoot() bool {
//	currentUser, err := user.Current()
//	if err != nil {
//		return false
//	}
//	return currentUser.Uid == "0"
//}

func init() {
	//if !isRoot() {
	//	fmt.Println("ERROR: Use sudo run this command like \"sudo simple-file-server\"")
	//	os.Exit(-1)
	//}
	//res, _ := module.ProcessList()
	//res, _ := module.ProcessConnectionList()
	//res, _ := module.ProcessSSHSessions()
	//console.GenerateSuccessData(res)
	//os.Exit(0)
}

func main() {
	log.Init()
	config.Init()
	if err := command.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
