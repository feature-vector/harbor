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

type executeResult struct {
	Code int64           `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (r executeResult) Error() string {
	return fmt.Sprintf("lark api execute error [code: %d] [msg: %s]", r.Code, r.Msg)
}

func executePostApi(ctx context.Context, url string, body map[string]interface{}) (*executeResult, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	tk, err := FetchAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http error code: %d body: %s", resp.StatusCode, string(respBody)))
	}
	r := &executeResult{}
	err = json.Unmarshal(respBody, r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, r
	}
	return r, nil
}

func executeGetApi(ctx context.Context, url string) (*executeResult, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	tk, err := FetchAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http error code: %d body: %s", resp.StatusCode, string(respBody)))
	}
	r := &executeResult{}
	err = json.Unmarshal(respBody, r)
	if err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, r
	}
	return r, nil
}
