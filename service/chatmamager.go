package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ChatManager struct {
	mu          sync.RWMutex
	Port        int
	userChatMap map[uuid.UUID]*ChatGLM
	canRun      bool
}

var (
	instance *ChatManager
)

const (
	EXPIRES_TIME_MINUTES = 30
)

func init() {
	instance = newChatManager()
}

func newChatManager() *ChatManager {
	r := &ChatManager{
		Port:        8000,
		userChatMap: make(map[uuid.UUID]*ChatGLM),
		canRun:      true,
	}

	go func() {
		timerHandler(r)
	}()
	return r
}

func GetChatManager() *ChatManager {
	return instance
}

func (cm *ChatManager) Stop() {
	cm.canRun = false
}

func (cm *ChatManager) GetSession(id uuid.UUID) *ChatGLM {
	fmt.Println("GetSession: " + id.String())

	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if value, ok := cm.userChatMap[id]; ok {
		value.Timestamp = time.Now()
		return value
	}
	return nil
}

func (cm *ChatManager) NewSession() uuid.UUID {
	id := uuid.New()

	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.userChatMap[id] = NewChatGLM()

	return id
}

func timerHandler(cm *ChatManager) {
	// 遍历所有 session
	for cm.canRun {
		cm.mu.RLock()
		for key, value := range cm.userChatMap {
			// 判断 session 是否过期
			now := time.Now()
			minutes := now.Sub(value.Timestamp).Minutes()
			if minutes > EXPIRES_TIME_MINUTES {
				// 删除过期 session
				fmt.Println("即将删除过期会话：" + key.String())
				cm.mu.RUnlock()
				cm.mu.Lock()
				delete(cm.userChatMap, key)
				cm.mu.Unlock()
				cm.mu.RLock()
			}
		}
		cm.mu.RUnlock()
		time.Sleep(time.Second)
	}
}
