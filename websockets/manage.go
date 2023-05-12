package websockets

import (
	"github.com/gorilla/websocket"
	"sync"
	"sync/atomic"
	"time"
)

type HandlerFunc interface {
	OnOpen(UniqueCode string)
	OnMessage(UniqueCode string, message []byte)
	OnClose(UniqueCode string)
}

func InitWebSocketManager(handler HandlerFunc) *WebSocketManager {
	wsManager := &WebSocketManager{Handler: handler}
	go wsManager.heartbeat()
	return wsManager
}

type WebSocketManager struct {
	count int64

	clientGroup sync.Map // map[string]*Client

	Handler HandlerFunc
}

func (m *WebSocketManager) onConnect(client *Client) {
	if client.uniqueCode == "" {
		client.message <- &messageData{data: nil, dataType: websocket.CloseMessage}
		return
	}
	atomic.AddInt64(&m.count, 1)
	m.safeAddClientClientGroup(client)
	m.Handler.OnOpen(client.uniqueCode)
}

func (m *WebSocketManager) onMessage(client *Client, message []byte) {
	m.Handler.OnMessage(client.uniqueCode, message)
}

func (m *WebSocketManager) disConnect(client *Client) {
	atomic.AddInt64(&m.count, -1)
	m.safeRemoveClientGroup(client.uniqueCode)
	m.Handler.OnClose(client.uniqueCode)
}

func (m *WebSocketManager) heartbeat() {
	for {
		m.clientGroup.Range(func(key, value interface{}) bool {
			value.(*Client).message <- &messageData{data: nil, dataType: websocket.PingMessage}
			return true
		})
		time.Sleep(15 * time.Second)
	}
}

func (m *WebSocketManager) safeRemoveClientGroup(UniqueCode string) {
	if client, ok := m.clientGroup.LoadAndDelete(UniqueCode); ok {
		close(client.(*Client).message)
	}
}

func (m *WebSocketManager) safeAddClientClientGroup(client *Client) {
	m.clientGroup.Store(client.uniqueCode, client)
}

func (m *WebSocketManager) registerClient(UniqueCode string, conn *websocket.Conn) {
	client := &Client{uniqueCode: UniqueCode, conn: conn, manager: m, message: make(chan *messageData)}
	go client.Read()
	go client.Write()
	m.onConnect(client)
}

// SendMessage 为指定客户端发送消息 UTF-8 编码
// UniqueCode 连接的唯一标识
// message 需要发送的消息
func (m *WebSocketManager) SendMessage(UniqueCode string, message []byte) {
	if client, ok := m.clientGroup.Load(UniqueCode); ok {
		client.(*Client).message <- &messageData{data: message, dataType: websocket.TextMessage}
	}
}

// BroadcastMessage 广播消息 UTF-8 编码
// message 需要广播的消息
// excluded 排除广播的连接的唯一标识
func (m *WebSocketManager) BroadcastMessage(message []byte, excluded ...string) {
	excludedMap := make(map[string]struct{}, len(excluded))
	for i := range excluded {
		excludedMap[excluded[i]] = struct{}{}
	}

	m.clientGroup.Range(func(key, value interface{}) bool {
		if _, ok := excludedMap[key.(string)]; !ok {
			value.(*Client).message <- &messageData{data: message, dataType: websocket.TextMessage}
		}
		return true
	})
}

// CloseWebsocketConn 关闭指定客户端得连接
// UniqueCode 连接的唯一标识
func (m *WebSocketManager) CloseWebsocketConn(UniqueCode string) {
	if client, ok := m.clientGroup.Load(UniqueCode); ok {
		client.(*Client).message <- &messageData{data: nil, dataType: websocket.CloseMessage}
	}
}

// SendBinary 为指定客户端发送二进制数据
// UniqueCode 连接的唯一标识
// binary 需要发送的二进制数据
func (m *WebSocketManager) SendBinary(UniqueCode string, binary []byte) {
	if client, ok := m.clientGroup.Load(UniqueCode); ok {
		client.(*Client).message <- &messageData{data: binary, dataType: websocket.BinaryMessage}
	}
}

// BroadcastBinary 广播二进制数据
// message 需要广播的二进制数据
// excluded 排除广播的连接的唯一标识
func (m *WebSocketManager) BroadcastBinary(binary []byte, excluded ...string) {
	excludedMap := make(map[string]struct{}, len(excluded))
	for i := range excluded {
		excludedMap[excluded[i]] = struct{}{}
	}

	m.clientGroup.Range(func(key, value interface{}) bool {
		if _, ok := excludedMap[key.(string)]; !ok {
			value.(*Client).message <- &messageData{data: binary, dataType: websocket.BinaryMessage}
		}
		return true
	})
}

// GetWSCount 获取连接数
func (m *WebSocketManager) GetWSCount() int64 {
	return atomic.LoadInt64(&m.count)
}
