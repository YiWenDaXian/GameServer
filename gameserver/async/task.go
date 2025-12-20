package async

import (
	"gameserver/common"
	"sync"
)

type TaskQueue[T common.TaskId] struct {
	ch chan T
}

func (tq *TaskQueue[T]) Select(task T) {
	select {
	case tq.ch <- task:
	default:
		// 队列满时的策略（示例：丢弃/统计/阻塞）
		tq.ch <- task // 或者记录并丢弃，按需调整
	}
}

type CoreTaskMap[T common.TaskId] struct {
	mu          sync.Mutex
	CoreMap     map[string]*TaskQueue[T]
	ProcessTask func(T)
}

func (ctm *CoreTaskMap[T]) INIT() {
	ctm.CoreMap = make(map[string]*TaskQueue[T])
}

func (ctm *CoreTaskMap[T]) worker(q *TaskQueue[T]) {
	for msg := range q.ch {
		ctm.ProcessTask(msg)
	}
}
func (ctm *CoreTaskMap[T]) GenChannel(taskId string) *TaskQueue[T] {
	q, ok := ctm.CoreMap[taskId]
	if !ok {
		ctm.mu.Lock()
		q = &TaskQueue[T]{ch: make(chan T, 1024)}
		ctm.CoreMap[taskId] = q
		go ctm.worker(q)
		ctm.mu.Unlock()
	}
	return q
}

//func (core *CoreTaskMap[T]) Sum(nums []T) T {
//	var total T
//	for _, num := range nums {
//		total += num
//	}
//	return total
//}
