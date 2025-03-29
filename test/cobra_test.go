package test

import (
	"fmt"
	"github.com/spf13/cobra"
	"gojks/cmd"
	"testing"
)

func TestCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmd.AddConfig)
	rootCmd.AddCommand(cmd.LsConfig)
	rootCmd.AddCommand(cmd.DelConfig)
	rootCmd.AddCommand(cmd.UseConfig)
	rootCmd.AddCommand(cmd.UseLs)
	rootCmd.AddCommand(cmd.Publish)

	// 模拟命令行输入
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"add", "http://localhost:8500", "admin:admin"})
	//rootCmd.SetArgs([]string{"del", "3"})
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
