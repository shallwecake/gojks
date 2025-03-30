package test

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojks/operation"
	"testing"
)

func TestCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(operation.AddConfig)
	rootCmd.AddCommand(operation.LsConfig)
	rootCmd.AddCommand(operation.DelConfig)
	rootCmd.AddCommand(operation.UseConfig)
	rootCmd.AddCommand(operation.UseLs)
	rootCmd.AddCommand(operation.Publish)

	// 模拟命令行输入
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"add", "http://localhost:8500", "admin:admin"})
	rootCmd.SetArgs([]string{"del", "3"})
	//rootCmd.SetArgs([]string{"ls"})
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"s"})
	//rootCmd.SetArgs([]string{"use", "4"})
	//rootCmd.SetArgs([]string{"uls"})
	//执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
