package ifunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	WebhookURL = "https://open.feishu.cn/open-apis/bot/v2/hook/a9dd87a0-ddb8-4224-a8b1-057ba580d72a"
)

// SendMessageToFeishu 发送消息到飞书机器人
func SendMessageToFeishu(webhookURL string, message interface{}) error {
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
