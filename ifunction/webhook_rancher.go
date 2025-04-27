package ifunction

import (
	"log"
	"time"
)

func RancherWebhook(has bool, name string, whk string) {
	var state string
	if has {
		state = "成功"
	} else {
		state = "失败"
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	// 构造卡片消息
	var message LarkCardMessage
	message.MsgType = "interactive"

	// 设置卡片配置
	message.Card.Config.WideScreenMode = true

	// 设置卡片标题
	message.Card.Header.Title.Tag = "plain_text"
	message.Card.Header.Title.Content = "rancher容器"

	// 添加卡片内容
	message.Card.Elements = append(message.Card.Elements, map[string]interface{}{
		"tag": "div",
		"text": map[string]interface{}{
			"content": "【" + name + "】启动" + state + "\t" + timeStr,
			"tag":     "plain_text",
		},
	})

	if err := sendMessageToFeishu(whk, message); err != nil {
		log.Fatal("webhook通知失败", err)
	}
}
