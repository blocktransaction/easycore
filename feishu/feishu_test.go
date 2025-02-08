package feishu

import "testing"

func TestFeishu(t *testing.T) {
	hookUrl := "https://open.feishu.cn/open-apis/bot/v2/hook/2e822cc0-4bda-4bd8-98c3-1fd3c2aa7ec6"
	message := NewMessage(hookUrl, MethodPost).EnableTime().SetPostTitle("B2Broker品种监控")
	message.SetContent(message.SetPostMsgContents([]string{"你好1", "你好2", "你好3"})).SendMessage()

}
