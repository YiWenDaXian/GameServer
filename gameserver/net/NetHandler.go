package net

import (
	"encoding/json"
	"fmt"

	"gameserver/common"

	"github.com/cloudwego/netpoll"
)

// Handler 实现了netpoll事件处理接口
type Handler struct {
	dispatcher Dispatcher
}
type ConnHandle struct {
	conn    netpoll.Connection
	writeCh chan []byte
}

var conns map[string]*ConnHandle

func (h *Handler) Init() {
	h.dispatcher.Init()
	conns = make(map[string]*ConnHandle)
}

// OnRead 处理读取事件
func (h *Handler) OnRead(conn netpoll.Connection) error {

	// 读取数据（只读取当前可用字节，避免等待固定大小导致阻塞）
	reader := conn.Reader()
	ln := reader.Len()
	if ln == 0 {
		// 没有可用数据，尽快返回，避免阻塞事件循环
		return nil
	}
	buffer, _ := reader.Peek(ln)
	data := make([]byte, len(buffer))
	copy(data, buffer)
	skipErr := reader.Skip(len(buffer))
	if skipErr != nil {
		conn.Close()
	}

	fmt.Printf("Received: %s\n", string(data))
	var rq common.Request
	err1 := json.Unmarshal(data, &rq)
	if err1 != nil {
		return err1
	}
	uid := rq.GenTaskId()
	if conns[uid] == nil {
		writeQueue := make(chan []byte, 100)
		conns[uid] = &ConnHandle{
			conn:    conn,
			writeCh: writeQueue,
		}
		go conns[uid].corePush()
	}
	// 分发请求到对应的用户队列
	// 发送响应
	// 发送响应 - 使用正确的方法名
	h.dispatcher.Dispatch(rq)
	return nil
}

// OnHup 处理连接关闭事件
func (h *Handler) OnHup(conn netpoll.Connection) error {
	fmt.Println("Connection closed")
	return nil
}
func Push(uid string, msg []byte) {
	conn, exit := conns[uid]
	if !exit {
		fmt.Println("Connection !exit")
		return
	}
	conn.writeCh <- msg

}
func (c *ConnHandle) corePush() {
	for data := range c.writeCh {
		_, writeErr := c.conn.Writer().WriteBinary(data)
		if writeErr != nil {
			fmt.Println("writeErr  ")
		}
		flushErr := c.conn.Writer().Flush()
		if flushErr != nil {
			fmt.Println("flushErr  ")
		}
	}
}
