package main

import (
	"fmt"
	"github.com/shallwecake/gojks/ifunction"
)

func main() {
	// 添加子命令到根命令
	ifunction.RootCmd.AddCommand(ifunction.AddConfig)
	ifunction.RootCmd.AddCommand(ifunction.DeleteConfig)
	ifunction.RootCmd.AddCommand(ifunction.LsConfig)
	ifunction.RootCmd.AddCommand(ifunction.PublishApp)
	ifunction.RootCmd.AddCommand(ifunction.PublishApps)
	// 执行根命令
	if err := ifunction.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
