package server

import (
	"bytes"
	"cqhttp-client/pkg/logutil"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type SendRawMessage struct {
	MessageType string `json:"message_type,omitempty"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id,omitempty"`
	Message     string `json:"message"`
	AutoEscape  bool   `json:"auto_escape,omitempty"`
}

type SendMessage struct {
	MessageType string  `json:"message_type,omitempty"`
	UserId      int64   `json:"user_id"`
	GroupId     int64   `json:"group_id,omitempty"`
	Message     Message `json:"message"`
	AutoEscape  bool    `json:"auto_escape,omitempty"`
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,string"`
}

func sendMessage(address string, toId int64, msgType string, message string) error {
	msg := Message{
		Type: msgType,
		Data: map[string]interface{}{
			msgType: message,
		},
	}
	sendMsg := SendMessage{
		UserId:  toId,
		Message: msg,
	}
	buf, _ := json.Marshal(sendMsg)
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/send_msg", address), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	tmp, _ := httputil.DumpRequest(req, true)
	logutil.Debug(string(tmp))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		logutil.Warn("Post send_msg statusCode:", resp.StatusCode)
	}

	return nil
}

func rawSendMessage(address string, body io.Reader) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/send_msg", address), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	tmp, _ := httputil.DumpRequest(req, true)
	logutil.Debug(string(tmp))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		logutil.Warn("Post send_msg statusCode:", resp.StatusCode)
	}

	return nil
}
