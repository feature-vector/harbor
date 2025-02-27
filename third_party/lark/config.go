package lark

var (
	larkAppId     string
	larkAppSecret string
)

func Init(appId string, appSecret string) {
	larkAppId = appId
	larkAppSecret = appSecret
}
