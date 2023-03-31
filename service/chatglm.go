package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/manucorporat/sse"
)

var (
	lockChatGLM sync.RWMutex
)

const (
	API_URL        = "http://127.0.0.1:8000"
	STREAM_API_URL = "http://127.0.0.1:8000/stream"
)

type ChatGLM struct {
	History   [][]string
	Timestamp time.Time
}

// ChatGLM 请求参数
type ApiPromptRequest struct {
	Prompt  string     `json:"prompt"`
	History [][]string `json:"history"`
}

// ChatGLM 响应结果
type ApiPromptResponse struct {
	Response string     `json:"response"`
	History  [][]string `json:"history"`
	Status   int        `json:"status"`
	Time     string     `json:"time"`
}

type ApiStreamPromptResponse struct {
	Message string `json:message`
}

func NewChatGLM() *ChatGLM {
	return &ChatGLM{
		History:   [][]string{},
		Timestamp: time.Now(),
	}
}

func handleErr(err error, prefix string) string {
	fmt.Println(prefix, err)
	return err.Error()
}

func (glm *ChatGLM) createRequestJson(prompt string) ([]byte, error) {
	request := ApiPromptRequest{
		Prompt:  prompt,
		History: glm.History,
	}
	jsonData, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (glm *ChatGLM) Prompt(prompt string) string {
	lockChatGLM.Lock()
	defer lockChatGLM.Unlock()

	//构造请求
	jsonData, err := glm.createRequestJson(prompt)

	if err != nil {
		return handleErr(err, "创建 JSON 失败：")
	}

	//发送请求
	resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return handleErr(err, "请求 ChatGLM 时发生错误：")
	}
	defer resp.Body.Close()

	//读取响应
	if resp.StatusCode != http.StatusOK {
		errMsg := "请求 ChatGLM 时发生错误：" + resp.Status
		fmt.Println(errMsg)
		return errMsg
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleErr(err, "读取 ChatGLM 响应时发生错误：")
	}

	//解析响应
	var response ApiPromptResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return handleErr(err, "解析 ChatGLM 接口响应内容时出错：")
	}

	glm.History = response.History
	return response.Response
}

func (glm *ChatGLM) StremPrompt(prompt string, outputStream io.Writer) {
	lockChatGLM.Lock()
	defer lockChatGLM.Unlock()

	//构造请求
	jsonData, err := glm.createRequestJson(prompt)

	if err != nil {
		handleErr(err, "创建 JSON 失败：")
		return
	}

	//发送请求
	resp, err := http.Post(STREAM_API_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		handleErr(err, "请求 ChatGLM 时发生错误：")
		return
	}
	defer resp.Body.Close()

	//读取响应
	if resp.StatusCode != http.StatusOK {
		errMsg := "请求 ChatGLM 时发生错误：" + resp.Status
		fmt.Println(errMsg)
	}

	//将 ChatGLM 的输出拷贝两份，一份输出给客户端，一份作为历史记录保存
	historyBuff := new(bytes.Buffer)

	flusher, _ := outputStream.(http.Flusher)
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err == nil {
			slice := buffer[:n]
			historyBuff.Write(slice)
			outputStream.Write(slice)
			flusher.Flush()
		} else {
			break
		}
	}

	events, err := sse.Decode(historyBuff)
	if err == nil && len(events) > 0 {
		d := events[len(events)-1].Data
		response := getString(d)
		var streamResponse ApiStreamPromptResponse
		json.Unmarshal([]byte(response), &streamResponse)
		glm.History = append(glm.History, []string{prompt, streamResponse.Message})
		fmt.Print("User: ")
		fmt.Println(prompt)
		fmt.Print("ChatGLM: ")
		fmt.Println(streamResponse.Message)
		return
	}

	fmt.Println("解析 Stream 出错，原因：" + err.Error())
}

func getString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
