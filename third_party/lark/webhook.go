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

func PostTextMessageToWebhook(ctx context.Context, webhook string, text string) error {
	body := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": text,
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
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
