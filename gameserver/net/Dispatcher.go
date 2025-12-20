package net

import (
	"fmt"
	"sync"
	"time"
)

type userQueue struct {
	ch    chan Request
	timer *time.Timer
}

type Dispatcher struct {
	mu     sync.Mutex
	queues map[string]*userQueue
}

func (d *Dispatcher) Init() {
	d.queues = make(map[string]*userQueue)
}

func (d *Dispatcher) Dispatch(msg Request) {
	d.mu.Lock()

	q, ok := d.queues[msg.UID]
	if !ok {
		q = &userQueue{ch: make(chan Request, 1024)}
		d.queues[msg.UID] = q
		go d.worker(msg.UID, q)
	}
	d.mu.Unlock()

	// 非阻塞入队，可根据策略阻塞/drop/扩容
	select {
	case q.ch <- msg:
	default:
		// 队列满时的策略（示例：丢弃/统计/阻塞）
		q.ch <- msg // 或者记录并丢弃，按需调整
	}
}

func (d *Dispatcher) worker(userID string, q *userQueue) {
	for msg := range q.ch {
		processMessage(msg)
	}
}
func processMessage(m Request) {
	// 在 worker goroutine 中安全地操作连接的写入
	// 处理业务...
	conn, ok := conns[m.UID]
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
