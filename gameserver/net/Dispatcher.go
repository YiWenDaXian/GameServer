package net

import (
	"fmt"
	"gameserver/async"
	"gameserver/common"
)

type Dispatcher struct {
	maps async.CoreTaskMap[common.Request]
}

func (d *Dispatcher) Init() {
	d.maps.INIT()
	d.maps.ProcessTask = processMessage
}

func (d *Dispatcher) Dispatch(msg common.Request) {

	// 非阻塞入队，可根据策略阻塞/drop/扩容
	q := d.maps.GenChannel(msg.GenTaskId())
	q.Select(msg)
}

func processMessage(m common.Request) {
	// 在 worker goroutine 中安全地操作连接的写入
	// 处理业务...
	conn, ok := conns[m.GenTaskId()]
	if !ok || conn == nil {
		fmt.Printf("No connection for uid %s, skip requestId: %d\n", m.UID, m.RequestId)
		return
	}
	writer := conn.Writer()
	if writer == nil {
		fmt.Printf("No writer for uid %s, skip requestId: %d\n", m.UID, m.RequestId)
		return
	}
	fmt.Printf("Processed uid %s requestId: %d\n", m.UID, m.RequestId)
	writer.WriteBinary([]byte("Processed requestId: " + fmt.Sprintf("%d", m.RequestId)))
	writer.Flush()
}
