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
	fmt.Printf("Processed uid %s requestId: %d\n", m.UID, m.RequestId)
	re := []byte(fmt.Sprintf("Processed requestId: %d", m.RequestId))
	Push(m.UID, re)
}
