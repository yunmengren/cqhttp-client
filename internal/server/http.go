package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"cqhttp-client/pkg/logutil"

	"github.com/buger/jsonparser"
)

type Handle struct {
	client        *http.Client
	CqhttpAddress string
	keywords      []string
	toID          int64
}

func (h *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logutil.Debugf("[%s] ", request.Method)
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, request.Body)
	data := buffer.Bytes()

	postType, err := jsonparser.GetString(data, "post_type")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	switch postType {
	case "message":
		err = h.handleMessage(data)
		if err != nil {
			logutil.Error("handleMessage error:", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	case "meta_event":
		err = h.handleMetaEvent(data)
		if err != nil {
			logutil.Error("handleMetaEvent error:", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		logutil.Warn("not support post_type:", postType)
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *Handle) handleMessage(data []byte) error {
	rawMessage, err := jsonparser.GetString(data, "raw_message")
	if err != nil {
		return err
	}
	for _, keyword := range h.keywords {
		if ok := strings.Contains(rawMessage, keyword); ok {
			logutil.Warn("出现关键消息:", keyword)
			sendMsg := SendRawMessage{
				UserId:  h.toID,
				Message: rawMessage,
			}
			buf, err := json.Marshal(sendMsg)
			if err != nil {
				return err
			}
			body := bytes.NewReader(buf)
			req, err := http.NewRequest("POST", "http://"+h.CqhttpAddress+"/send_msg", body)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			//tmp, _ := httputil.DumpRequest(req, true)
			//fmt.Println(string(tmp))
			resp, err := h.client.Do(req)
			if err != nil {
				return err
			} else if resp.StatusCode != http.StatusOK {
				logutil.Warn("Post send_msg statusCode:", resp.StatusCode)
			}
			out := bytes.NewBuffer(nil)
			io.Copy(out, resp.Body)
			resp.Body.Close()
			logutil.Debug("response body:", out.String())
			break
		}
	}

	return nil
}

func (h *Handle) handleMetaEvent(data []byte) error {
	metaEventType, err := jsonparser.GetString(data, "meta_event_type")
	if err != nil {
		return err
	}
	if metaEventType == "heartbeat" {
		logutil.Debug("recieve heartbeat from cqhttp")
	} else {
		logutil.Warn("not support meta_event_type:", metaEventType)
	}
	return nil
}
