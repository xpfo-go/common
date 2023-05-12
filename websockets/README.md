# Websockets

## 一 介绍

### 一个简单方便易用并发安全得 websocket 组件，对 github.com/gorilla/websocket 包进行得二次封装，生命周期与 js 对应。

## 二 使用

### 1. 实现 WsHandlerFunc

#### 实现 WsHandlerFunc接口，编写业务逻辑。 接口实现得方法为 OnOpen，建立连接时回调，OnMessage 接收到消息时回调，OnClose 连接关闭时回调。

``` go
var WsHandler *WsHandlerFunc

func init() {
    WsHandler = &WsHandlerFunc{}
}

type WsHandlerFunc struct {}
func (ws *WsHandlerFunc) OnOpen(UniqueCode string) {}
func (ws *WsHandlerFunc) OnMessage(UniqueCode string, message []byte) {}
func (ws *WsHandlerFunc) OnClose(UniqueCode string) {}
```

### 2. 实现 WebSocketManager

#### 实现 WebSocketManager 对websocket连接进行管理，快捷方便得使用内置函数

``` go
var WsManage *websockets.WebSocketManager

func init() {
    WsManage = websockets.InitWebSocketManager(WsHandler)
}
```


### 3. 升级协议

### 无论使用哪种 web 框架均可使用此方法升级协议 （ UpGrande ）

#### 对于原生 net / http包，使用 UpGrande 方法升级协议，将 http / https 协议升级为 ws / wss 协。uniqueCode 为此连接得唯一标识，推荐使用幂等得 uniqueCode（比如user_id），将连接交给 WebSocketManager 管理。

``` go
http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
    websockets.UpGrande(w, r, uniqueCode, WsManage)
})
```

### 为方便使用，以下为 go 常用 web 框架协议升级方法

#### 对于 Gin，使用 UpGrandeForGin 方法升级协议。

``` go
func Foo(c *gin.Context) {
    websockets.UpGrandeForGin(c, uniqueCode, WsManage)
}
```

#### 对于 Beego，使用 UpGrandeForBeego 方法升级协议。

``` go
import "github.com/beego/beego/v2/server/web/context"

func Foo(ctx *context.Context) {
    websockets.UpGrandeForBeego(ctx, uniqueCode, WsManage)
}
```

#### 对于 GoFrame，使用 UpGrandeForGoFrame 方法升级协议。

``` go
func Foo(r *ghttp.Request) {
    websockets.UpGrade(r.Response.ResponseWriter, r.Request, uniqueCode, WsManage)
}
```

### 4. WebSocketManager 内置函数

* 为指定连接发送消息 (UniqueCode)， message 为 utf-8 编码 []byte。

``` go
func (m *WebSocketManager) SendMessage(UniqueCode string, message []byte)
``` 

* 为所有连接发送消息，message 为 utf-8 编码 []byte，excluded 则是排除发送得连接得 UniqueCode

``` go
func (m *WebSocketManager) BroadcastMessage(message []byte, excluded ...string) {}
```

* 为指定连接发送二进制数据

``` go
func (m *WebSocketManager) SendBinary(UniqueCode string, binary []byte) {}
```

* 为所有连接发送二进制数据

``` go
func (m *WebSocketManager) BroadcastBinary(binary []byte, excluded ...string) {}
```

* 断开指定连接

``` go
func (m *WebSocketManager) CloseWebsocketConn(UniqueCode string) {}
```

* 获取连接数

``` go
func (m *WebSocketManager) GetWSCount() int64 {}
```