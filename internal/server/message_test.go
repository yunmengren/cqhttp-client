package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestSendMessage(t *testing.T) {
	msg := Message{
		Type: "json",
		Data: map[string]interface{}{"data": "{\"app\":\"com.tencent.structmsg\",\"config\":{\"autosize\":true,\"ctime\":1646636423,\"forward\":true,\"token\":\"6e8da4fcc2adc441ec3addec3bb43c9c\",\"type\":\"normal\"},\"desc\":\"新闻\",\"extra\":{\"app_type\":1,\"appid\":100951776,\"uin\":1622912909},\"meta\":{\"news\":{\"action\":\"\",\"android_pkg_name\":\"\",\"app_type\":1,\"appid\":100951776,\"ctime\":1646636423,\"desc\":\"UP主：夏白熊烨-房间号：45378\",\"jumpUrl\":\"https://b23.tv/K3mjNd8\",\"preview\":\"https://pic.ugcimg.cn/e9047036331c8b361435a17a3f64e224/jpg1\",\"source_icon\":\"\",\"source_url\":\"\",\"tag\":\"哔哩哔哩\",\"title\":\"白熊受苦受难中\",\"uin\":1622912909}},\"prompt\":\"[分享]白熊受苦受难中\",\"ver\":\"0.0.0.1\",\"view\":\"news\"}"},
	}
	sendMsg := SendMessage{
		UserId:  1215697684,
		Message: msg,
	}
	buf, _ := json.Marshal(sendMsg)
	body := bytes.NewBuffer(buf)
	err := rawSendMessage("127.0.0.1:5700", body)
	if err != nil {
		fmt.Println("sendMessage error:", err)
	}
}
