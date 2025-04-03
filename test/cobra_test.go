package test

import (
	"fmt"
	"github.com/shallwecake/gojks/ifunction"
	"github.com/spf13/cobra"
	"testing"
)

func TestCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(ifunction.AddConfig)
	rootCmd.AddCommand(ifunction.PublishApp)
	rootCmd.AddCommand(ifunction.PublishApps)

	// 模拟命令行输入
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"add", "http://localhost:8500", "admin:admin"})
	rootCmd.SetArgs([]string{"add", "whk", "http://localhost:8500"})
	//rootCmd.SetArgs([]string{"del", "3"})
	//rootCmd.SetArgs([]string{"ls"})
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"s"})
	//rootCmd.SetArgs([]string{"use", "1"})
	//rootCmd.SetArgs([]string{"uls"})
	//执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
