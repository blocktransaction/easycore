package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	MethodPost = "post"
	MethodText = "text"
)

type Message struct {
	//webhook地址
	webhookurl string
	//消息类型
	msgType string
	//消息内容
	content interface{}
	// 标题
	title string
	//是否显示时间
	enableTime bool
	//当前时间
	currDatetime string
	//默认标题
	defaultTitle string
}

func NewMessage(webhookurl, msgType string) *Message {
	return &Message{
		webhookurl:   webhookurl,
		msgType:      msgType,
		currDatetime: time.Now().Format(time.DateTime),
		defaultTitle: "消息提醒",
	}
}

// 设置默认消息标题
func (m *Message) SetDefaultTitle(defaultTitle string) *Message {
	m.defaultTitle = defaultTitle
	return m
}

// 设置hook链接
func (m *Message) SetHookUrl(webhookurl string) *Message {
	m.webhookurl = webhookurl
	return m
}

// 设置消息类型
func (m *Message) SetMgType(msgType string) *Message {
	m.msgType = msgType
	return m
}

// 设置消息标题
func (m *Message) SetPostTitle(title string) *Message {
	if title == "" {
		title = m.defaultTitle
		if m.enableTime {
			title = m.currDatetime + "" + title
		}
	}
	m.title = title
	if m.enableTime {
		m.title = m.currDatetime + " " + m.title
	}
	return m
}

// 设置开启显示时间
func (m *Message) EnableTime() *Message {
	m.enableTime = true
	return m
}

// 设置post方法的消息内容
func (m *Message) SetPostMsgContents(contents []string) []Content {
	msgContents := make([]Content, 0)
	for _, content := range contents {
		msgContents = append(msgContents, Content{
			Tag:      MethodText,
			UnEscape: false,
			Text:     content + "\r\n",
		})
	}
	return msgContents
}

// 设置最终要发送的消息内容(包括post和text方法)
func (m *Message) SetContent(content interface{}) *Message {
	switch m.msgType {
	case MethodPost:
		msg := PostMsg{
			MsgType: m.msgType,
		}
		msg.Content.Post.ZhCn.Title = m.title
		msg.Content.Post.ZhCn.Content = content.([]Content)
		for k, v := range msg.Content.Post.ZhCn.Content {
			if v.Tag == "" {
				msg.Content.Post.ZhCn.Content[k].Tag = MethodText
			}
		}
		content = msg
	case MethodText:
		msg := TextMsg{
			MsgType: m.msgType,
		}
		msg.Content.Text = content.(string)
		if m.enableTime {
			msg.Content.Text = m.currDatetime + " " + content.(string)
		}
		content = msg
	}
	m.content = content
	return m
}

// 发送飞书消息
func (m *Message) SendMessage() (bool, error) {
	var result msgResult

	bs, err := json.Marshal(m.content)
	if err != nil {
		return false, err
	}
	if m.msgType == MethodPost {
		temp := string(bs)
		temp = strings.ReplaceAll(temp, "\"content\":[", "\"content\":[[")
		temp = strings.ReplaceAll(temp, "]", "]]")
		bs = []byte(temp)
	}

	req, err := http.NewRequest(http.MethodPost, m.webhookurl, bytes.NewBuffer(bs))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}
	if result.Code != 0 {
		return false, fmt.Errorf("send message failed, code: %d, msg: %s", result.Code, result.Msg)
	}
	return true, nil
}
