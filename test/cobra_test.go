package test

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"gojks"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(createConfig)
	rootCmd.AddCommand(lsConfig)
	rootCmd.AddCommand(delConfig)
	rootCmd.AddCommand(publish)
	rootCmd.AddCommand(search)
	rootCmd.AddCommand(useConfig)
	rootCmd.AddCommand(useLs)

	// 模拟命令行输入
	rootCmd.SetArgs([]string{"pub", "test-jenkins-Pipeline"})

	//rootCmd.SetArgs([]string{"create", "test", "https://jenkins.gw-greenenergy.com", "pangwangbin:wongbin123"})
	//rootCmd.SetArgs([]string{"create", "pre", "https://jenkins.gw-greenenergy.com", "pangwangbin:wongbin123"})
	//rootCmd.SetArgs([]string{"create", "local", "http://localhost:8500", "admin:admin"})

	//rootCmd.SetArgs([]string{"del", "3"})
	//rootCmd.SetArgs([]string{"ls"})
	//rootCmd.SetArgs([]string{"pub"})
	//rootCmd.SetArgs([]string{"s", "test"})
	//rootCmd.SetArgs([]string{"use", "4"})
	//rootCmd.SetArgs([]string{"uls"})
	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}

// 定义子命令
var publish = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"pub"}, // 定义别名
	Short:   "发布应用",
	Args:    cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		engine := main.InitDb()
		defer main.CloseDb(engine)
		id := main.GetUse(engine)
		config := main.Get(engine, id)
		auth := &main.Auth{
			Username: config.Username,
			ApiToken: config.Password,
		}

		jenkins := main.NewJenkins(auth, config.Url)
		names, _ := jenkins.Query(appName)

		jenkinsURL := "http://localhost:8500"
		jobName := "test-jenkins-Pipeline"
		username := "admin"
		apiToken := "admin"
		ctx := context.Background()
		jenkins := gojenkins.CreateJenkins(nil, jenkinsURL, username, apiToken)
		_, err := jenkins.Init(ctx)
		if err != nil {
			panic("连接 Jenkins 失败: " + err.Error())
		}
		fmt.Println("Jenkins 连接成功")
	},
}

var search = &cobra.Command{
	Use:     "search",
	Aliases: []string{"q", "s"}, // 定义别名
	Short:   "查询任务",
	Args:    cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		engine := main.InitDb()
		defer main.CloseDb(engine)
		id := main.GetUse(engine)
		config := main.Get(engine, id)

		auth := &main.Auth{
			Username: config.Username,
			ApiToken: config.Password,
		}

		jenkins := main.NewJenkins(auth, config.Url)
		names, _ := jenkins.Query(name)

		for _, name := range names {
			fmt.Println(name)
		}
	},
}

var createConfig = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "cre"},
	Short:   "创建配置",
	Args:    cobra.ExactArgs(3), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		env := args[0]
		url := args[1]
		auth := args[2]
		slice := strings.Split(auth, ":")
		info := &main.Config{
			Env:      env,
			Url:      url,
			Username: slice[0],
			Password: slice[1],
		}
		engine := main.InitDb()
		defer main.CloseDb(engine)
		main.Save(engine, info)
	},
}

var lsConfig = &cobra.Command{
	Use:   "ls",
	Short: "遍历配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := main.InitDb()
		defer main.CloseDb(engine)
		main.Ls(engine)
	},
}

var delConfig = &cobra.Command{
	Use:   "del",
	Short: "删除配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := main.InitDb()
		defer main.CloseDb(engine)
		main.Del(engine, args[0])
	},
}

var useConfig = &cobra.Command{
	Use:   "use",
	Short: "使用配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := main.InitDb()
		defer main.CloseDb(engine)
		main.Use(engine, args[0])
	},
}

var useLs = &cobra.Command{
	Use:   "uls",
	Short: "使用配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := main.InitDb()
		defer main.CloseDb(engine)
		main.UseLs(engine)
	},
}
