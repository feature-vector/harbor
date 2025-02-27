package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/feature-vector/harbor/base/hc"
	"io/ioutil"
	"net/http"
)

type Client struct {
	AppId     string
	AppSecret string
}

type Event struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

func NewEvent(name string) *Event {
	return &Event{Name: name, Params: map[string]interface{}{}}
}

func (e *Event) WithParam(k string, v interface{}) *Event {
	e.Params[k] = v
	return e
}

func (c *Client) PostEvents(ctx context.Context, appInstanceId string, userId string, events []*Event) error {
	body := map[string]interface{}{
		"app_instance_id": appInstanceId,
		"user_id":         userId,
		"events":          events,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	u := fmt.Sprintf(`https://www.google-analytics.com/mp/collect?firebase_app_id=%s&api_secret=%s`, c.AppId, c.AppSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("firebase.PostEvents http error: %d %s", resp.StatusCode, string(respBody))
	}
	return nil
}
