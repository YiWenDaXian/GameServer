package main

import (
	"context"
	"fmt"
	"gameserver/net"

	"github.com/cloudwego/netpoll"
)

func main() {
	// 创建监听器
	shouldReturn := newFunction()
	if shouldReturn {
		return
	}
}

func newFunction() bool {
	listener, err := netpoll.CreateListener("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to create listener: %v\n", err)
		return true
	}
	defer listener.Close()

	// 创建事件处理器
	//handler := &net.Handler{}
	handler := new(net.Handler)

	// 创建事件循环 - 使用正确的API和参数
	looper, err := netpoll.NewEventLoop(func(ctx context.Context, conn netpoll.Connection) error {
		// 当有新连接时调用OnRead
		return handler.OnRead(conn)
	})
	if err != nil {
		fmt.Printf("Failed to create event loop: %v\n", err)
		return true
	}

	// 启动服务 - 使用正确的API和参数
	fmt.Println("Starting TCP server on :8080 (listening on all interfaces)")
	err = looper.Serve(listener)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
	return false
}
