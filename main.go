package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojks/storage"
	"log"
	"strings"
)

func main() {
	// 添加子命令到根命令
	//rootCmd.AddCommand(publish)
	//rootCmd.AddCommand(create)
	rootCmd.AddCommand(ls)

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

// 定义子命令
var publish = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"pub"}, // 定义别名
	Short:   "发布应用",
	Args:    cobra.ExactArgs(2), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		env := args[1]
		fmt.Printf("开始发布应用：%s\n", appName)
		fmt.Printf("目标环境：%s\n", env)
		// 在这里添加实际的发布逻辑
		fmt.Println("发布完成！")
	},
}

var create = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "cre"},
	Short:   "创建环境",
	Args:    cobra.ExactArgs(3), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		url := args[1]
		auth := args[2]
		slice := strings.Split(auth, ":")
		info := map[int]string{
			0: url,
			1: slice[0],
			2: slice[1],
		}
		if storage.Save(env, info) {
			fmt.Println("创建成功")
		} else {
			log.Print("创建失败")
		}
	},
}

var ls = &cobra.Command{
	Use:   "ls",
	Short: "遍历环境",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		storage.Ls()
	},
}

var del = &cobra.Command{
	Use:   "del",
	Short: "删除环境",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		storage.Del(args[0])
	},
}
