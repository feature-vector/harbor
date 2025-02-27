package lark

import (
	"context"
)

type CardElement struct {
	Tag  string   `json:"tag"`
	Text CardText `json:"text"`
}

type CardText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type CardHeader struct {
	Template string    `json:"template"`
	Title    CardTitle `json:"title"`
}

type CardTitle struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Card struct {
	Elements []CardElement `json:"elements"`
	Header   CardHeader    `json:"header"`
}

type CardCallbackAction struct {
	Value map[string]interface{} `json:"value"`
	Tag   string                 `json:"tag"`
}

type CardCallback struct {
	OpenId        string             `json:"open_id"`
	UserId        string             `json:"user_id"`
	OpenMessageId string             `json:"open_message_id"`
	OpenChatId    string             `json:"open_chat_id"`
	TenantKey     string             `json:"tenant_key"`
	Token         string             `json:"token"`
	Action        CardCallbackAction `json:"action"`
}

func UpdateCard(ctx context.Context, token string, card *Card) error {
	url := "https://open.feishu.cn/open-apis/interactive/v1/card/update"
	body := map[string]interface{}{
		"token": token,
		"card":  card,
	}
	_, err := executePostApi(ctx, url, body)
	return err
}
