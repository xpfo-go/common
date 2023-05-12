package websockets

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrande = websocket.Upgrader{
	//设置允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpGrade(w http.ResponseWriter, r *http.Request, UniqueCode string, ws *WebSocketManager) error {
	conn, err := upGrande.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	ws.registerClient(UniqueCode, conn)
	return nil
}

func UpGrandeForGin(c *gin.Context, UniqueCode string, ws *WebSocketManager) error {
	//创建连接
	conn, err := upGrande.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	ws.registerClient(UniqueCode, conn)
	return nil
}

func UpGrandeForBeego(c *context.Context, UniqueCode string, ws *WebSocketManager) error {
	//创建连接
	conn, err := upGrande.Upgrade(c.ResponseWriter, c.Request, nil)
	if err != nil {
		return err
	}
	ws.registerClient(UniqueCode, conn)
	return nil
}

func UpGrandeForGoFrame(c *ghttp.Request, UniqueCode string, ws *WebSocketManager) error {
	//创建连接
	conn, err := upGrande.Upgrade(c.Response.ResponseWriter, c.Request, nil)
	if err != nil {
		return err
	}
	ws.registerClient(UniqueCode, conn)
	return nil
}
