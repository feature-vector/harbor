package lark

import (
	"context"
	"fmt"
)

func FetchTableFields(ctx context.Context, appToken string, tableId string) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/fields", appToken, tableId)
	_, err := executeGetApi(ctx, url)
	return err
}

func CreateRecord(ctx context.Context, appToken string, tableId string, fields map[string]interface{}) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records", appToken, tableId)
	body := map[string]interface{}{
		"fields": fields,
	}
	_, err := executePostApi(ctx, url, body)
	return err
}
