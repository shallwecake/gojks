package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shallwecake/gojks/ifunction"
	"net/http"
	"testing"
)

func TestWhk(t *testing.T) {
	engine := ifunction.InitDb()
	defer ifunction.CloseDbEngine(engine)
	url := ifunction.GetConf(engine, ifunction.Web_Hook).Url
	ifunction.RancherWebhook(true, "go-test", url)
}

func TestWhk01(t *testing.T) {

	// 定义飞书卡片消息结构
	type LarkCardMessage struct {
		MsgType string `json:"msg_type"`
		Card    struct {
			Config struct {
				WideScreenMode bool `json:"wide_screen_mode"`
			} `json:"config"`
			Header struct {
				Title struct {
					Tag     string `json:"tag"`
					Content string `json:"content"`
				} `json:"title"`
			} `json:"header"`
			Elements []interface{} `json:"elements"`
		} `json:"card"`
	}

	// 替换为你的飞书 Webhook URL
	webhookURL := ""

	// 构造卡片消息
	var message LarkCardMessage
	message.MsgType = "interactive"

	// 设置卡片配置
	message.Card.Config.WideScreenMode = true

	// 设置卡片标题
	message.Card.Header.Title.Tag = "plain_text"
	message.Card.Header.Title.Content = "这是一张卡片消息"

	// 添加卡片内容
	message.Card.Elements = append(message.Card.Elements, map[string]interface{}{
		"tag": "div",
		"text": map[string]interface{}{
			"content": "这是卡片中的内容部分。",
			"tag":     "plain_text",
		},
	})

	// 将消息序列化为 JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("HTTP 请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusOK {
		fmt.Println("卡片消息发送成功！")
	} else {
		fmt.Printf("卡片消息发送失败，状态码: %d\n", resp.StatusCode)
	}

}
