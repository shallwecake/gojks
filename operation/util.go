package operation

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"gojks/jenkins_suggest"
	"gojks/store"
	"os"
	"strconv"
	"sync"
	"time"
)

func isJobInQueue(jenkins *gojenkins.Jenkins, jobName string, ctx context.Context, i int) {
	// 获取队列对象
	queue, _ := jenkins.GetQueue(ctx)
	// 遍历队列中的所有任务
	items := queue.Raw.Items

	is := true

	if len(items) == 0 {
		is = false
	}

	for _, item := range items {
		if item.Task.Name == jobName {
			//fmt.Printf("Task ID: %d\n", item.ID)
			//fmt.Printf("Task Name: %s\n", item.Task.Name)
			//fmt.Printf("Why in Queue: %s\n", item.Why)
			//fmt.Printf("Stuck: %t\n", item.Stuck)
			//fmt.Printf("URL: %s\n", item.URL)
			//fmt.Println("-------------------------------")
			is = true
		}

	}

	if is {
		time.Sleep(1 * time.Second)
		timePrint("队列中", i)
		isJobInQueue(jenkins, jobName, ctx, i+1)
	}
}

func timePrint(name string, i int) {
	// 打印计时信息，\r 将光标移回行首
	fmt.Printf("\r%s: %02d 秒", name, i)
	// 刷新输出缓冲区，确保内容立即显示
	fmt.Print("\033[?25l") // 隐藏光标（可选）
}

func msgPrint(name string) {
	// 打印计时信息，\r 将光标移回行首
	fmt.Printf("\r%s", name)
	// 刷新输出缓冲区，确保内容立即显示
	fmt.Print("\033[?25l") // 隐藏光标（可选）
}

func clearPrint() {
	fmt.Print("\r\033[K")
}

func Suggest(config *store.Config, name string) []string {
	auth := &jenkins_suggest.Auth{
		Username: config.Username,
		ApiToken: config.Password,
	}
	jenkins := jenkins_suggest.NewJenkins(auth, config.Url)
	names, _ := jenkins.Query(name)
	if len(names) == 0 {
		return []string{}
	}
	return names
}

func SyncPublish(jenkins *gojenkins.Jenkins, ctx context.Context, suggest []string, wg *sync.WaitGroup) {
	fmt.Printf("序号\t名称\n")
	for i, item := range suggest {
		fmt.Printf("%d\t%s\n", i, item)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("请输入构建的序号：")
	if scanner.Scan() { // 读取一行
		input := scanner.Text() // 获取文本（自动去除换行符）
		i, _ := strconv.Atoi(input)
		name := suggest[i]

		// 触发指定任务（Job）的构建
		_, err := jenkins.BuildJob(ctx, name, nil)
		if err != nil {
			panic("触发构建失败: " + err.Error())
		}

		msgPrint("正在准备构建,请稍等...")
		time.Sleep(500 * time.Millisecond) // 避免CPU跑满
		clearPrint()
		//清空整行
		isJobInQueue(jenkins, name, ctx, 0)
		go func() {
			running := true
			var loopCount int
			job, _ := jenkins.GetJob(ctx, name)
			lastBuild, _ := job.GetLastBuild(ctx)
			for running {
				loopCount++
				number := lastBuild.GetBuildNumber()
				build, _ := job.GetBuild(ctx, number)
				if build.IsRunning(ctx) {
					timePrint("构建中", loopCount)
				} else {
					clearPrint()
					fmt.Println("构建成功")
					fmt.Print("\033[?25h") // 开启光标
					running = false
				}
				time.Sleep(1 * time.Second) // 避免CPU跑满
			}
			wg.Done() // 协程结束时计数器-1
		}()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取错误:", err)
	}
}
