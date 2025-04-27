package ifunction

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"sync"
	"time"
)

// RootCmd 定义根命令
var RootCmd = &cobra.Command{
	Use:   "jenkins",
	Short: "Jenkins命令行工具",
}

// PublishApp Publish 定义子命令
var PublishApp = &cobra.Command{
	Use: "pub",
	//Aliases: []string{"pub"}, // 定义别名
	Short: "发布单个应用",
	Args:  cobra.ExactArgs(1), // 确保必须提供1个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		engine := InitDb()
		defer CloseDbEngine(engine)
		config := GetConf(engine, Jen_kins)
		if len(config.Url) == 0 {
			log.Printf("请配置%s", Jen_kins)
		}
		suggest := Suggest(config, name)
		if len(suggest) == 0 {
			// 重试
			suggest = Suggest(config, name)
		}

		if len(suggest) == 0 {
			fmt.Printf("没有找到构建名称 %s\n", name)
		} else {
			ctx := context.Background()
			jenkins := gojenkins.CreateJenkins(nil, config.Url, config.Username, config.Password)
			_, err := jenkins.Init(ctx)
			if err != nil {
				panic("连接 Jenkins 失败: " + err.Error())
			}
			// 构建
			PublishJob(jenkins, ctx, suggest)
		}
	},
}

var PublishApps = &cobra.Command{
	Use: "pubs",
	//Aliases: []string{"pub"}, // 定义别名
	Short: "发布多个应用,输入全名并用英文逗号分隔",
	Args:  cobra.ExactArgs(1), // 不限制参数
	Run: func(cmd *cobra.Command, args []string) {
		names := strings.Split(args[0], ",")
		engine := InitDb()
		defer CloseDbEngine(engine)
		config := GetConf(engine, Jen_kins)
		ctx := context.Background()
		jenkins := gojenkins.CreateJenkins(nil, config.Url, config.Username, config.Password)
		_, err := jenkins.Init(ctx)
		if err != nil {
			panic("连接 Jenkins 失败: " + err.Error())
		}

		var wg sync.WaitGroup

		for _, name := range names {
			time.Sleep(100 * time.Millisecond)
			wg.Add(1)
			go func() {

				_, err := jenkins.BuildJob(ctx, name, nil)

				if err != nil {
					fmt.Printf("构建失败：%s\n", name)
				}

				wg.Done()
			}()

		}

		wg.Wait()

	},
}

var AddConfig = &cobra.Command{
	Use: "add",
	//Aliases: []string{"new", "cre"},
	Short: "创建配置",
	//Args:  cobra.ExactArgs(3),  确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		// 参数数量校验
		if len(args) < 2 || len(args) > 3 {
			fmt.Printf("参数数量错误：需要 2 到 3 个参数，实际接收到 %d 个", len(args))
		}
		category := args[0]
		var url string
		if len(args) > 1 {
			url = args[1]
		} else {
			url = ""
		}
		var slice []string
		if len(args) > 2 {
			slice = strings.Split(args[2], ":")
		} else {
			slice = []string{"", ""}
		}

		info := &Conf{
			Type:     category,
			Url:      url,
			Username: slice[0],
			Password: slice[1],
		}

		engine := InitDb()
		defer CloseDbEngine(engine)
		SaveOrUpdate(engine, info)
	},
}

var LsConfig = &cobra.Command{
	Use:   "ls",
	Short: "查看配置",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		engine := InitDb()
		defer CloseDbEngine(engine)
		ListConf(engine)
	},
}

var DeleteConfig = &cobra.Command{
	Use:   "del",
	Short: "删除配置",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		engine := InitDb()
		defer CloseDbEngine(engine)
		DelConf(engine, args[0])
	},
}
