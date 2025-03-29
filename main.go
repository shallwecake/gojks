package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"gojks/jenkins_suggest"
	"gojks/store"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 添加子命令到根命令
	rootCmd.AddCommand(createConfig)
	rootCmd.AddCommand(lsConfig)
	rootCmd.AddCommand(delConfig)
	rootCmd.AddCommand(publish)
	rootCmd.AddCommand(useConfig)
	rootCmd.AddCommand(useLs)

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
	Args:    cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		engine := store.InitDb()
		defer store.CloseDb(engine)
		id := store.GetUse(engine)
		config := store.Get(engine, id)
		auth := &jenkins_suggest.Auth{
			Username: config.Username,
			ApiToken: config.Password,
		}

		suggest := jenkins_suggest.NewJenkins(auth, config.Url)
		names, _ := suggest.Query(name)

		if len(names) == 0 {

			fmt.Printf("没有找到%s...", name)

		} else {
			fmt.Printf("序号\t名称\n")
			for i, item := range names {
				fmt.Printf("%d\t%s\n", i, item)
			}

			scanner := bufio.NewScanner(os.Stdin)

			fmt.Print("请输入构建的序号：")
			if scanner.Scan() { // 读取一行
				input := scanner.Text() // 获取文本（自动去除换行符）
				i, _ := strconv.Atoi(input)
				jname := names[i]

				ctx := context.Background()
				jenkins := gojenkins.CreateJenkins(nil, config.Url, config.Username, config.Password)
				_, err := jenkins.Init(ctx)
				if err != nil {
					panic("连接 Jenkins 失败: " + err.Error())
				}
				fmt.Println("Jenkins 连接成功")

				// 触发指定任务（Job）的构建
				_, err = jenkins.BuildJob(ctx, jname, nil)
				if err != nil {
					panic("触发构建失败: " + err.Error())
				}

				fmt.Printf("正在构建中，请稍后...")

				go func() {
					running := true

					for running {
						job, _ := jenkins.GetJob(ctx, jname)
						lastBuild, _ := job.GetLastBuild(ctx)

						if !lastBuild.IsRunning(ctx) {
							result := lastBuild.GetResult()
							fmt.Printf("构建%s", result)
							running = false
						}

						time.Sleep(1 * time.Second) // 避免CPU跑满
					}

				}()
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("读取错误:", err)
			}
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
		info := &store.Config{
			Env:      env,
			Url:      url,
			Username: slice[0],
			Password: slice[1],
		}
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Save(engine, info)
	},
}

var lsConfig = &cobra.Command{
	Use:   "ls",
	Short: "遍历配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Ls(engine)
	},
}

var delConfig = &cobra.Command{
	Use:   "del",
	Short: "删除配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Del(engine, args[0])
	},
}

var useConfig = &cobra.Command{
	Use:   "use",
	Short: "使用配置",
	Args:  cobra.ExactArgs(1), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.Use(engine, args[0])
	},
}

var useLs = &cobra.Command{
	Use:   "uls",
	Short: "使用配置",
	Args:  cobra.ExactArgs(0), // 确保必须提供两个参数
	Run: func(cmd *cobra.Command, args []string) {
		engine := store.InitDb()
		defer store.CloseDb(engine)
		store.UseLs(engine)
	},
}
