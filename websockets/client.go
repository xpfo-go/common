package websockets

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	uniqueCode string
	conn       *websocket.Conn
	manager    *WebSocketManager

	message chan *messageData
}

type messageData struct {
	dataType int
	data     []byte
}

func (c *Client) Read() {
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			c.manager.disConnect(c)
			break
		}

		c.manager.onMessage(c, msg)
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.message:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(msg.dataType, msg.data); err != nil {
				return
			}
		}
	}

}
