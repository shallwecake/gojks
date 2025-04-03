package ifunction

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"os"
	"strconv"
	"strings"
	"time"
)

func isJobInQueue(jenkins *gojenkins.Jenkins, jobName string, ctx context.Context, second int) {
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
		for i := 1; i <= 3; i++ {
			second += 1
			time.Sleep(1 * time.Second)
			timeMsgPrint("队列中", second)
		}
		isJobInQueue(jenkins, jobName, ctx, second)
	}
}

func timeMsgPrint(name string, i int) {
	fmt.Print("\r\033[K")
	// 打印计时信息，\r 将光标移回行首
	fmt.Printf("\r%s: %02d 秒", name, i)
	// 刷新输出缓冲区，确保内容立即显示
	fmt.Print("\033[?25l") // 隐藏光标（可选）
}

func msgPrint(name string) {
	fmt.Print("\r\033[K")
	// 打印计时信息，\r 将光标移回行首
	fmt.Printf("\r%s", name)
	// 刷新输出缓冲区，确保内容立即显示
	fmt.Print("\033[?25l") // 隐藏光标（可选）
}

func refreshPrint() {
	fmt.Print("\r\033[K")
}

func closePrint() {
	fmt.Print("\033[?25h") // 开启光标
}

func Suggest(config *Conf, name string) []string {
	auth := &Auth{
		Username: config.Username,
		ApiToken: config.Password,
	}
	jenkins := NewJenkins(auth, config.Url)
	names, _ := jenkins.Query(name)
	if len(names) == 0 {
		return []string{}
	}
	return names
}

func PublishJob(jenkins *gojenkins.Jenkins, ctx context.Context, suggest []string) {

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
		time.Sleep(2 * time.Second) // 避免CPU跑满
		isJobInQueue(jenkins, name, ctx, 0)
		MonitorJenkins(jenkins, ctx, name)
		MonitorRancher(name)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取错误:", err)
	}

	closePrint()
}

func MonitorRancher(name string) {
	conf := GetConf(Engine, Ran_Cher)
	access := &Access{
		AccessKey: conf.Username,
		SecretKey: conf.Password,
	}

	rancher := NewRancher(access, conf.Url)

	switch {
	case strings.Contains(name, "test"):
		clusters := rancher.GetClusters("test")
		data := clusters.Data[0]
		pods := rancher.GetPodList(data.ID, data.Name).Data

		for _, pod := range pods {
			pname := pod.Metadata.Labels["app"]
			//id := pod.ID
			phase := pod.Status.Phase
			startTime := pod.Status.StartTime
			fmt.Println(pname + " " + phase + " " + startTime)
		}

	case strings.Contains(name, "pre"):

	default:

	}
	//// 判断是否包含 "test" 或 "pre"
	//containsTest := strings.Contains(name, "test")
	//containsPre := strings.Contains(name, "pre")
}

func MonitorJenkins(jenkins *gojenkins.Jenkins, ctx context.Context, name string) {
	running := true
	var loopCount int
	job, _ := jenkins.GetJob(ctx, name)
	lastBuild, _ := job.GetLastBuild(ctx)
	for running {
		number := lastBuild.GetBuildNumber()
		build, _ := job.GetBuild(ctx, number)
		if build.IsRunning(ctx) {
			for i := 1; i <= 5; i++ {
				loopCount += 1
				time.Sleep(1 * time.Second)
				timeMsgPrint("构建中", loopCount)
			}
		} else {
			refreshPrint()
			fmt.Printf("构建%s", switchResult(build.GetResult()))
			isSendMsg(build)
			running = false
		}
	}
}

func isSendMsg(build *gojenkins.Build) {
	has := false
	if build.GetResult() == "SUCCESS" {
		has = true
	}
	Webhook(has, build.Job)
}

func switchResult(item interface{}) (result string) {
	switch item {
	case "FAILURE":
		result = "失败"
	case "SUCCESS":
		result = "成功"
	}
	return
}

func nowTime() string {
	// 时区
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// 获取当前时间
	currentTime := time.Now().In(loc)

	// 格式化时间为 yyyy-dd-mm hh:mm:ss
	return currentTime.Format("2006-02-01 15:04:05")
}
