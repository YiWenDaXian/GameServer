package net

import (
	"encoding/json"
	"fmt"

	"github.com/cloudwego/netpoll"
)

// Handler 实现了netpoll事件处理接口
type Handler struct {
	dispatcher Dispatcher
}

var conns map[string]netpoll.Connection

type Request struct {
	uid       string                 `json:"uid"`
	requestId int                    `json:"requestId"`
	parma     map[string]interface{} `json:"parma"`
}

// OnRead 处理读取事件
func (h *Handler) OnRead(conn netpoll.Connection) error {
	if conns == nil {
		conns = make(map[string]netpoll.Connection)
	}

	// 读取数据
	buffer, _ := conn.Reader().Peek(4096)
	data := make([]byte, len(buffer))
	copy(data, buffer)
	conn.Reader().Skip(len(buffer))

	fmt.Printf("Received: %s\n", string(data))
	var rq Request
	err1 := json.Unmarshal([]byte(data), &rq)
	if err1 != nil {
		return err1
	}
	uid := rq.uid
	if conns[uid] == nil {
		conns[uid] = conn
	}
	// 分发请求到对应的用户队列
	// 发送响应
	// 发送响应 - 使用正确的方法名
	h.dispatcher.Dispatch(rq)
	writer := conn.Writer()
	_, err := writer.WriteBinary(data)
	if err != nil {
		return err
	}
	// 刷新缓冲区
	writer.Flush()

	return nil
}

// OnHup 处理连接关闭事件
func (h *Handler) OnHup(conn netpoll.Connection) error {
	fmt.Println("Connection closed")
	return nil
}
