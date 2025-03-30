package operation

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"gojks/store"
	"strings"
	"sync"
)

// RootCmd 定义根命令
var RootCmd = &cobra.Command{
	Use:   "jenkins",
	Short: "Jenkins命令行工具",
}

// Publish 定义子命令
var Publish = &cobra.Command{
	Use: "pub",
	//Aliases: []string{"pub"}, // 定义别名
	Short: "发布应用",
	Args:  cobra.ExactArgs(1), // 确保必须提供1个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		engine := store.InitDb()
		defer store.CloseDb(engine)
		id := store.GetUse(engine)
		config := store.Get(engine, id)

		suggest := Suggest(config, name)

		var wg sync.WaitGroup
		wg.Add(1) // 计数器+1

		if len(suggest) == 0 {

			fmt.Printf("没有找到%s...\n", name)

		} else {
			ctx := context.Background()
			jenkins := gojenkins.CreateJenkins(nil, config.Url, config.Username, config.Password)
			_, err := jenkins.Init(ctx)
			if err != nil {
				panic("连接 Jenkins 失败: " + err.Error())
			}
			fmt.Println("Jenkins 连接成功")

			// 构建
			SyncPublish(jenkins, ctx, suggest, &wg)
		}
		wg.Wait()
	},
}

var AddConfig = &cobra.Command{
	Use: "add",
	//Aliases: []string{"new", "cre"},
	Short: "创建配置",
	Args:  cobra.ExactArgs(2), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		auth := args[1]
		slice := strings.Split(auth, ":")
		info := &store.Config{
			Url:      url,
			Username: slice[0],
			Password: slice[1],
		}
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Save(engine, info)
	},
}

var LsConfig = &cobra.Command{
	Use:   "ls",
	Short: "遍历配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Ls(engine)
	},
}

var DelConfig = &cobra.Command{
	Use:   "del",
	Short: "删除配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Del(engine, args[0])
	},
}

var UseConfig = &cobra.Command{
	Use:   "use",
	Short: "使用配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Use(engine, args[0])
	},
}

var UseLs = &cobra.Command{
	Use:   "uls",
	Short: "使用配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.UseLs(engine)
	},
}
