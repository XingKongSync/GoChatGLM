package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Prompt  string `json:"prompt"`
	Session string `json:"Session"`
}

func main() {
	request := Request{
		Prompt:  "你好",
		Session: "ffc36263-d5cd-4f7e-9ca3-64f8a4ba4daa",
	}
	jsonData, _ := json.Marshal(request)

	resp, err := http.Post("http://localhost:3001/api/streamchat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err == nil {
			slice := buffer[:n]
			fmt.Println(string(slice))
		} else {
			break
		}
	}
}
