package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	// 添加子命令到根命令
	//rootCmd.AddCommand(publish)
	//rootCmd.AddCommand(create)
	//rootCmd.AddCommand(ls)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// 定义根命令
var rootCmd = &cobra.Command{
	Use:   "jenkins",
	Short: "Jenkins命令行工具",
}
