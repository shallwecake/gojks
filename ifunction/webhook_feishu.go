package ifunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"
	"io"
	"log"
	"net/http"
)

func Webhook(has bool, job *gojenkins.Job) {
	webhookURL := GetConf(Engine, Web_Hook).Url
	if len(webhookURL) > 0 {
		return
	}
	commitInfos := []CommitInfo{
		{
			Title:       "提交信息",
			Description: "feat: 根据电站容量限制申报电量",
			Date:        "2023-04-02 17:26:27",
		},
		{
			Title:       "Merge branch '250319_新疆日报'",
			Description: "",
			Date:        "2023-04-02 17:26:54",
		},
	}

	message := Message{
		MsgType: "post",
		Post: PostContent{
			Title: "api-declare 代码推送 test",
			Content: []Paragraph{
				{
					Type: "header",
					Text: commitInfos[0].Title,
				},
				{
					Type: "text",
					Text: fmt.Sprintf("#### %s\n%s", commitInfos[0].Description, commitInfos[0].Date),
				},
				{
					Type: "header",
					Text: commitInfos[1].Title,
				},
				{
					Type: "text",
					Text: fmt.Sprintf("%s", commitInfos[1].Date),
				},
			},
		},
	}
	if err := sendMessageToFeishu(webhookURL, message); err != nil {
		log.Fatal("webhook通知失败", err)
	}
}

// SendMessageToFeishu 发送消息到飞书机器人
func sendMessageToFeishu(webhookURL string, message interface{}) error {
	// 将消息体序列化为 JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message. Status code: %d", resp.StatusCode)
	}

	//fmt.Println("Message sent successfully!")
	return nil
}

type CommitInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type Message struct {
	MsgType string      `json:"msg_type"`
	Content interface{} `json:"content"`
	Text    string      `json:"text"`
	Post    PostContent `json:"post"`
}

type PostContent struct {
	Title   string      `json:"title"`
	Content []Paragraph `json:"content"`
}

type Paragraph struct {
	Type     string    `json:"tag"`
	Text     string    `json:"text"`
	Elements []Element `json:"elements"`
}

type Element struct {
	Type     string    `json:"tag"`
	Text     string    `json:"text"`
	Elements []Element `json:"elements"`
}
