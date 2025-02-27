package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/feature-vector/harbor/base/hc"
	"io/ioutil"
	"net/http"
)

const (
	EventTypeReceiveMessage = "im.message.receive_v1"
)

type Event struct {
	Schema string `json:"schema"`
	Header struct {
		EventId    string `json:"event_id"`
		EventType  string `json:"event_type"`
		CreateTime string `json:"create_time"`
		Token      string `json:"token"`
		AppId      string `json:"app_id"`
		TenantKey  string `json:"tenant_key"`
	} `json:"header"`
	Event struct {
		Sender struct {
			SenderId struct {
				UnionId string `json:"union_id"`
				UserId  string `json:"user_id"`
				OpenId  string `json:"open_id"`
			} `json:"sender_id"`
			SenderType string `json:"sender_type"`
			TenantKey  string `json:"tenant_key"`
		} `json:"sender"`
		Message struct {
			MessageId   string `json:"message_id"`
			RootId      string `json:"root_id"`
			ParentId    string `json:"parent_id"`
			CreateTime  string `json:"create_time"`
			ChatId      string `json:"chat_id"`
			ChatType    string `json:"chat_type"`
			MessageType string `json:"message_type"`
			Content     string `json:"content"`
			Mentions    []struct {
				Key string `json:"key"`
				Id  struct {
					UnionId string `json:"union_id"`
					UserId  string `json:"user_id"`
					OpenId  string `json:"open_id"`
				} `json:"id"`
				Name      string `json:"name"`
				TenantKey string `json:"tenant_key"`
			} `json:"mentions"`
		} `json:"message"`
	} `json:"event"`
}

func ReplyTextMessage(ctx context.Context, messageId string, text string) error {
	contentBytes, err := json.Marshal(map[string]interface{}{
		"text": text,
	})
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"msg_type": "text",
		"content":  string(contentBytes),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s/reply", messageId)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	tk, err := FetchAccessToken(ctx)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("http error code: %d body: %s", resp.StatusCode, string(respBody)))
	}
	return nil
}

func ReplyTemplateMessage(ctx context.Context, messageId string, templateId string, templateVariables map[string]interface{}) error {
	contentBytes, err := json.Marshal(map[string]interface{}{
		"type": "template",
		"data": map[string]interface{}{
			"template_id":       templateId,
			"template_variable": templateVariables,
		},
	})
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"msg_type": "interactive",
		"content":  string(contentBytes),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s/reply", messageId)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	tk, err := FetchAccessToken(ctx)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("http error code: %d body: %s", resp.StatusCode, string(respBody)))
	}
	return nil
}
