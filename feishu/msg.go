package feishu

type TextMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type PostMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string    `json:"title"`
				Content []Content `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

type Content struct {
	Tag      string `json:"tag"`
	Text     string `json:"text"`
	UnEscape bool   `json:"un_escape"`
}

type msgResult struct {
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
	Code          int         `json:"code"`
	Data          interface{} `json:"data"`
	Msg           string      `json:"msg"`
}
