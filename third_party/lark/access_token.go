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
	"time"
)

var (
	tenantAccessToken          = ""
	tenantAccessTokenExpiredAt = time.Now()
)

type accessTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func FetchAccessToken(ctx context.Context) (string, error) {
	if time.Now().Before(tenantAccessTokenExpiredAt) {
		return tenantAccessToken, nil
	}
	appId := larkAppId
	appSecret := larkAppSecret
	body := map[string]interface{}{
		"app_id":     appId,
		"app_secret": appSecret,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := hc.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("http error code: %d body: %s", resp.StatusCode, string(respBody)))
	} else {
		r := &accessTokenResp{}
		err = json.Unmarshal(respBody, r)
		if err != nil {
			return "", err
		}
		if r.Code != 0 {
			fmt.Println(string(respBody))
			return "", errors.New(fmt.Sprintf("lark error code: %d, msg: %s", r.Code, r.Msg))
		}
		tenantAccessToken = r.TenantAccessToken
		tenantAccessTokenExpiredAt = time.Now().Add(time.Duration(r.Expire) * time.Second)
		return tenantAccessToken, nil
	}
}
