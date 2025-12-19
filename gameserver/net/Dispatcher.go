package net

import (
	"fmt"
	"sync"
	"time"
)

type Msg struct {
	userID    string
	requestId string
	parma     map[string]interface{}
}

type userQueue struct {
	ch    chan Msg
	timer *time.Timer
}

type Dispatcher struct {
	mu     sync.Mutex
	queues map[string]*userQueue
	idle   time.Duration
}

func (d *Dispatcher) Dispatch(userID string, requestId string, parma map[string]interface{}) {
	d.mu.Lock()
	q, ok := d.queues[userID]
	if !ok {
		q = &userQueue{ch: make(chan Msg, 1024)}
		d.queues[userID] = q
		go d.worker(userID, q)
	}
	d.mu.Unlock()
	m := Msg{
		userID:    userID,
		requestId: requestId,
		parma:     parma,
	}
	// 非阻塞入队，可根据策略阻塞/drop/扩容
	select {
	case q.ch <- m:
	default:
		// 队列满时的策略（示例：丢弃/统计/阻塞）
		q.ch <- m // 或者记录并丢弃，按需调整
	}
}

func (d *Dispatcher) worker(userID string, q *userQueue) {
	idleTimer := time.NewTimer(d.idle)
	defer idleTimer.Stop()

	for {
		select {
		case msg := <-q.ch:
			// 收到消息则重置回收计时
			if !idleTimer.Stop() {
				select {
				case <-idleTimer.C:
				default:
				}
			}
			idleTimer.Reset(d.idle)

			// 这里顺序处理消息（业务逻辑）：
			processMessage(msg)

		case <-idleTimer.C:
			// 空闲超时，回收队列
			d.mu.Lock()
			delete(d.queues, userID)
			d.mu.Unlock()
			return
		}
	}
}

func processMessage(m Msg) {
	// 在 worker goroutine 中安全地操作连接的写入
	// 处理业务...
	conn := conns[m.userID]
	writer := conn.Writer()
	fmt.Printf(m.requestId)
	_, _ = writer.WriteBinary([]byte("Processed requestId: " + m.requestId))
	writer.Flush()
}
