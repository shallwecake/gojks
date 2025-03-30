package main

import (
	"fmt"
	"github.com/shallwecake/gojks/operation"
)

func main() {
	// 添加子命令到根命令
	operation.RootCmd.AddCommand(operation.AddConfig)
	operation.RootCmd.AddCommand(operation.LsConfig)
	operation.RootCmd.AddCommand(operation.DelConfig)
	operation.RootCmd.AddCommand(operation.UseConfig)
	operation.RootCmd.AddCommand(operation.UseLs)
	operation.RootCmd.AddCommand(operation.Publish)
	operation.RootCmd.AddCommand(operation.PublishAll)
	// 执行根命令
	if err := operation.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
