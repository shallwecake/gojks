package ifunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
