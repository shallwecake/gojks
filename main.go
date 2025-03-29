package main

import (
	"fmt"
	"gojks/cmd"
)

func main() {
	// 添加子命令到根命令
	cmd.RootCmd.AddCommand(cmd.CreateConfig)
	cmd.RootCmd.AddCommand(cmd.LsConfig)
	cmd.RootCmd.AddCommand(cmd.DelConfig)
	cmd.RootCmd.AddCommand(cmd.UseConfig)
	cmd.RootCmd.AddCommand(cmd.UseLs)
	cmd.RootCmd.AddCommand(cmd.Publish)
	// 执行根命令
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
