package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func NewServer(wwwPath *string) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)
	n.Use(negroni.NewStatic(http.Dir(*wwwPath)))
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/api/chat", chatHandler(formatter)).Methods("POST")
	mx.HandleFunc("/api/streamchat", streamChatHandler).Methods("GET")
	mx.HandleFunc("/api/session", sessionHandler(formatter)).Methods("GET")
}

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ChatRequest struct {
	Session string `json:"session"`
	Prompt  string `json:"prompt"`
}

type ChatResponse struct {
	CommonResponse
	Response string `json:"response"`
}

type SessionResponse struct {
	CommonResponse
	Session string `json:"session"`
}

func chatHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//解析请求
		var chatRequest ChatRequest
		body, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, &chatRequest)
		}
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, CommonResponse{Code: -1, Message: "解析请求时出错"})
			return
		}

		//处理请求
		id, err := uuid.Parse(chatRequest.Session)
		var chatGLM *ChatGLM = nil
		if err == nil {
			chatGLM = GetChatManager().GetSession(id)
		}
		if chatGLM == nil {
			formatter.JSON(w, http.StatusForbidden, CommonResponse{Code: -1, Message: "会话已过期"})
			return
		}
		response := chatGLM.Prompt(chatRequest.Prompt)
		formatter.JSON(w, http.StatusOK, ChatResponse{
			CommonResponse: CommonResponse{
				Code:    0,
				Message: "成功",
			},
			Response: response,
		})
	}
}

func streamChatHandler(w http.ResponseWriter, req *http.Request) {
	//解析请求
	var chatRequest ChatRequest
	// body, err := io.ReadAll(req.Body)
	// req.Body.Close()
	// if err == nil {
	// 	err = json.Unmarshal(body, &chatRequest)
	// }
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	args := req.URL.Query()
	chatRequest.Prompt = args["prompt"][0]
	chatRequest.Session = args["session"][0]

	//处理请求
	id, err := uuid.Parse(chatRequest.Session)
	var chatGLM *ChatGLM = nil
	if err == nil {
		chatGLM = GetChatManager().GetSession(id)
	}
	if chatGLM == nil {
		http.Error(w, "会话已过期", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	chatGLM.StremPrompt(chatRequest.Prompt, w)
}

func sessionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := GetChatManager().NewSession().String()

		formatter.JSON(w, http.StatusOK, SessionResponse{
			CommonResponse: CommonResponse{
				Code:    0,
				Message: "成功",
			},
			Session: id,
		})
	}
}
